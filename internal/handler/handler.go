package handler

import (
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"container-refresh/internal/config"
	"container-refresh/internal/docker"
	"container-refresh/internal/slack"
)

// UpdateRequest defines the expected JSON body for the /update endpoint.
type UpdateRequest struct {
	Token string `json:"token"`
}

// Handler holds dependencies for HTTP handlers.
type Handler struct {
	Config *config.Config
}

// NewHandler creates a new Handler with the given configuration.
func NewHandler(cfg *config.Config) *Handler {
	return &Handler{Config: cfg}
}

// UpdateHandler handles requests to the /update endpoint.
func (h *Handler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding JSON request: %v", err)
		http.Error(w, "Malformed JSON request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close() // Ensure the request body is closed

	if req.Token == "" {
		http.Error(w, "Missing token in request", http.StatusBadRequest)
		return
	}

	// Securely compare the received token with the configured token
	// subtle.ConstantTimeCompare returns 1 if the slices are equal, 0 otherwise.
	tokenMatch := subtle.ConstantTimeCompare([]byte(req.Token), []byte(h.Config.Token)) == 1

	if !tokenMatch {
		log.Printf("Unauthorized attempt: received token '%s', expected token starting with '%s...'", req.Token, h.Config.Token[:min(5, len(h.Config.Token))])
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	log.Println("Token validated. Starting container pull process...")

	// Docker client library is used now, executable name is no longer needed

	if err := docker.PullContainers(h.Config.Images); err != nil {
		log.Printf("Error pulling container images: %v", err)
		errorMsg := fmt.Sprintf("Failed to pull container images: %v", err)
		http.Error(w, errorMsg, http.StatusInternalServerError)
		
		// Send failure notification to Slack
		hostname, _ := os.Hostname()
		slackMsg := fmt.Sprintf("Containers not refreshed on %s: %s", hostname, errorMsg)
		if err := slack.SendMessage(slackMsg); err != nil {
			log.Printf("Failed to send Slack notification: %v", err)
		}
		return
	}
	log.Println("All container images pulled successfully.")

	log.Println("Starting container stop process...")
	if err := docker.StopContainers(h.Config.ContainerNames); err != nil {
		log.Printf("Error stopping containers: %v", err)
		errorMsg := fmt.Sprintf("Failed to stop containers: %v", err)
		http.Error(w, errorMsg, http.StatusInternalServerError)
		
		// Send failure notification to Slack
		hostname, _ := os.Hostname()
		slackMsg := fmt.Sprintf("Containers not refreshed on %s: %s", hostname, errorMsg)
		if err := slack.SendMessage(slackMsg); err != nil {
			log.Printf("Failed to send Slack notification: %v", err)
		}
		return
	}
	log.Println("All containers stopped successfully.")

	w.WriteHeader(http.StatusOK)
	fw, err := w.Write([]byte("Update process completed successfully."))
	if err != nil {
		log.Printf("Error writing response: %v, bytes written: %d", err, fw)
	}
	log.Println("Update process completed successfully and response sent.")

	// Send success notification to Slack
	hostname, _ := os.Hostname()
	slackMsg := fmt.Sprintf("Containers refreshed on %s", hostname)
	if err := slack.SendMessage(slackMsg); err != nil {
		log.Printf("Failed to send Slack notification: %v", err)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
