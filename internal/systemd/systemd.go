package systemd

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"container-refresh/internal/config"
)

// RestartServices iterates through a list of service configurations and restarts them.
func RestartServices(services []config.Service) error {
	if len(services) == 0 {
		log.Println("No services configured to restart.")
		return nil
	}

	var errorMessages []string
	log.Printf("Attempting to restart %d service(s).", len(services))
	for _, service := range services {
		// Use custom restart command if provided, otherwise default to systemctl restart
		restartCmd := fmt.Sprintf("systemctl restart %s", service.Name)
		if service.RestartCommand != "" {
			restartCmd = service.RestartCommand
		}

		log.Printf("Executing: %s", restartCmd)
		
		// Split the command into parts for exec.Command
		cmdParts := strings.Fields(restartCmd)
		if len(cmdParts) == 0 {
			errMsg := fmt.Sprintf("Invalid restart command for service '%s'", service.Name)
			log.Println(errMsg)
			errorMessages = append(errorMessages, errMsg)
			continue
		}

		cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
		output, err := cmd.CombinedOutput() // CombinedOutput includes both stdout and stderr
		if err != nil {
			errMsg := fmt.Sprintf("Failed to restart '%s': %v. Output: %s", service.Name, err, string(output))
			log.Println(errMsg)
			errorMessages = append(errorMessages, errMsg)
			// Continue to the next service even if one fails
			continue
		}
		
		// Restart command doesn't typically produce much output on success, but log if any
		if len(strings.TrimSpace(string(output))) > 0 {
			log.Printf("Output from restart command for %s: %s", service.Name, string(output))
		}
		log.Printf("Successfully initiated restart for '%s'.", service.Name)
	}

	if len(errorMessages) > 0 {
		return fmt.Errorf("encountered errors restarting one or more services:\n%s", strings.Join(errorMessages, "\n"))
	}
	log.Println("All configured services restarted (or attempted to restart) successfully.")
	return nil
}
