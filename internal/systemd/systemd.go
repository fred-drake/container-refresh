package systemd

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// RestartServices iterates through a list of systemd service names and restarts them.
func RestartServices(services []string) error {
	if len(services) == 0 {
		log.Println("No systemd services configured to restart.")
		return nil
	}

	var errorMessages []string
	log.Printf("Attempting to restart %d systemd service(s).", len(services))
	for _, serviceName := range services {
		log.Printf("Executing: systemctl restart %s", serviceName)
		cmd := exec.Command("systemctl", "restart", serviceName)
		output, err := cmd.CombinedOutput() // CombinedOutput includes both stdout and stderr
		if err != nil {
			errMsg := fmt.Sprintf("Failed to restart '%s': %v. Output: %s", serviceName, err, string(output))
			log.Println(errMsg)
			errorMessages = append(errorMessages, errMsg)
			// Continue to the next service even if one fails
			continue
		}
		// systemctl restart doesn't typically produce much output on success, but log if any
		if len(strings.TrimSpace(string(output))) > 0 {
			log.Printf("Output from systemctl restart %s: %s", serviceName, string(output))
		}
		log.Printf("Successfully initiated restart for '%s'.", serviceName)
	}

	if len(errorMessages) > 0 {
		return fmt.Errorf("encountered errors restarting one or more systemd services:\n%s", strings.Join(errorMessages, "\n"))
	}
	log.Println("All configured systemd services restarted (or attempted to restart) successfully.")
	return nil
}
