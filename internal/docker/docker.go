package docker

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"container-refresh/internal/config"
)

// PullContainers iterates through a list of container configurations and pulls them using the specified executable.
func PullContainers(executable string, containers []config.Container) error {
	if len(containers) == 0 {
		log.Println("No containers configured to pull.")
		return nil
	}

	var errorMessages []string
	log.Printf("Attempting to pull %d container(s) using '%s'.", len(containers), executable)
	for _, container := range containers {
		log.Printf("Executing: %s pull %s", executable, container.Image)
		cmd := exec.Command(executable, "pull", container.Image)
		output, err := cmd.CombinedOutput()
		if err != nil {
			errMsg := fmt.Sprintf("Failed to pull '%s' (name: %s) using '%s': %v. Output: %s", 
				container.Image, container.Name, executable, err, string(output))
			log.Println(errMsg)
			errorMessages = append(errorMessages, errMsg)
			// Continue to the next container even if one fails
			continue
		}
		log.Printf("Successfully pulled '%s' (name: %s) using '%s'. Output (last 2 lines):\n%s", 
			container.Image, container.Name, executable, getLastNLines(string(output), 2))
	}

	if len(errorMessages) > 0 {
		return fmt.Errorf("encountered errors pulling one or more containers using '%s':\n%s", 
			executable, strings.Join(errorMessages, "\n"))
	}
	log.Printf("All configured containers pulled (or attempted to pull) successfully using '%s'.", executable)
	return nil
}

// getLastNLines returns the last N lines of a string, useful for concise logging of command output.
func getLastNLines(s string, n int) string {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	if len(lines) <= n {
		return s
	}
	return strings.Join(lines[len(lines)-n:], "\n")
}
