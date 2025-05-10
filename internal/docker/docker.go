package docker

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// PullContainers iterates through a list of container image tags and pulls them using the specified executable.
func PullContainers(executable string, images []string) error {
	if len(images) == 0 {
		log.Println("No container images configured to pull.")
		return nil
	}

	var errorMessages []string
	log.Printf("Attempting to pull %d container image(s) using '%s'.", len(images), executable)
	for _, imageTag := range images {
		log.Printf("Executing: %s pull %s", executable, imageTag)
		cmd := exec.Command(executable, "pull", imageTag)
		output, err := cmd.CombinedOutput()
		if err != nil {
			errMsg := fmt.Sprintf("Failed to pull '%s' using '%s': %v. Output: %s", 
				imageTag, executable, err, string(output))
			log.Println(errMsg)
			errorMessages = append(errorMessages, errMsg)
			// Continue to the next container even if one fails
			continue
		}
		log.Printf("Successfully pulled '%s' using '%s'. Output (last 2 lines):\n%s", 
			imageTag, executable, getLastNLines(string(output), 2))
	}

	if len(errorMessages) > 0 {
		return fmt.Errorf("encountered errors pulling one or more container images using '%s':\n%s", 
			executable, strings.Join(errorMessages, "\n"))
	}
	log.Printf("All configured container images pulled (or attempted to pull) successfully using '%s'.", executable)
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

// StopContainers stops containers by their names using the specified executable.
func StopContainers(executable string, containerNames []string) error {
	if len(containerNames) == 0 {
		log.Println("No containers configured to stop.")
		return nil
	}

	var errorMessages []string
	log.Printf("Attempting to stop %d container(s) using '%s'.", len(containerNames), executable)
	for _, containerName := range containerNames {
		log.Printf("Executing: %s stop %s", executable, containerName)
		cmd := exec.Command(executable, "stop", containerName)
		output, err := cmd.CombinedOutput()
		if err != nil {
			errMsg := fmt.Sprintf("Failed to stop container '%s' using '%s': %v. Output: %s", 
				containerName, executable, err, string(output))
			log.Println(errMsg)
			errorMessages = append(errorMessages, errMsg)
			// Continue to the next container even if one fails
			continue
		}
		log.Printf("Successfully stopped container '%s' using '%s'. Output (last 2 lines):\n%s", 
			containerName, executable, getLastNLines(string(output), 2))
	}

	if len(errorMessages) > 0 {
		return fmt.Errorf("encountered errors stopping one or more containers using '%s':\n%s", 
			executable, strings.Join(errorMessages, "\n"))
	}
	log.Printf("All configured containers stopped (or attempted to stop) successfully using '%s'.", executable)
	return nil
}
