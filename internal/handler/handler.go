package handler

import (
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"container-refresh/internal/config"
	"container-refresh/internal/docker"
	"container-refresh/internal/systemd"
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

	containerExecutable := h.Config.Executable
	if containerExecutable == "" {
		containerExecutable = "docker"
	}

	if err := docker.PullContainers(containerExecutable, h.Config.Containers); err != nil {
		log.Printf("Error pulling containers: %v", err)
		http.Error(w, fmt.Sprintf("Failed to pull containers: %v", err), http.StatusInternalServerError)
		return
	}
	log.Println("All containers pulled successfully.")

	log.Println("Starting service restart process...")
	if err := systemd.RestartServices(h.Config.SystemdServices); err != nil {
		log.Printf("Error restarting services: %v", err)
		http.Error(w, fmt.Sprintf("Failed to restart services: %v", err), http.StatusInternalServerError)
		return
	}
	log.Println("All services restarted successfully.")

	w.WriteHeader(http.StatusOK)
	fw, err := w.Write([]byte("Update process completed successfully."))
	if err != nil {
		log.Printf("Error writing response: %v, bytes written: %d", err, fw)
	}
	log.Println("Update process completed successfully and response sent.")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
