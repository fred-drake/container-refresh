package docker

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/fsouza/go-dockerclient"
)

// PullContainers pulls multiple Docker images using the Docker client library
func PullContainers(images []string) error {
	if len(images) == 0 {
		log.Println("No container images configured to pull.")
		return nil
	}

	// Create a new Docker client
	client, err := docker.NewClientFromEnv()
	if err != nil {
		return fmt.Errorf("failed to create Docker client: %w", err)
	}

	var errorMessages []string
	log.Printf("Attempting to pull %d container image(s) using Docker client.", len(images))

	for _, imageTag := range images {
		log.Printf("Pulling image: %s", imageTag)

		// Pull the image
		err := client.PullImage(docker.PullImageOptions{
			Repository: imageTag,
		}, docker.AuthConfiguration{})

		if err != nil {
			errMsg := fmt.Sprintf("Failed to pull '%s': %v", imageTag, err)
			log.Println(errMsg)
			errorMessages = append(errorMessages, errMsg)
			// Continue to the next container even if one fails
			continue
		}

		log.Printf("Successfully pulled '%s'", imageTag)
	}

	if len(errorMessages) > 0 {
		return fmt.Errorf("encountered errors pulling one or more container images:\n%s",
			strings.Join(errorMessages, "\n"))
	}

	log.Println("All configured container images pulled (or attempted to pull) successfully.")
	return nil
}

// StopContainers stops multiple Docker containers by name using the Docker client library
func StopContainers(containerNames []string) error {
	if len(containerNames) == 0 {
		log.Println("No containers configured to stop.")
		return nil
	}

	// Create a new Docker client
	client, err := docker.NewClientFromEnv()
	if err != nil {
		return fmt.Errorf("failed to create Docker client: %w", err)
	}

	var errorMessages []string
	log.Printf("Attempting to stop %d container(s) using Docker client.", len(containerNames))

	// List all containers to find the ones we need to stop
	containers, err := client.ListContainers(docker.ListContainersOptions{All: true})
	if err != nil {
		return fmt.Errorf("failed to list containers: %w", err)
	}

	// Create a map of container names to IDs for quick lookup
	containerMap := make(map[string]string)
	for _, container := range containers {
		for _, name := range container.Names {
			// Container names from the API have a leading slash
			cleanName := strings.TrimPrefix(name, "/")
			containerMap[cleanName] = container.ID
		}
	}

	// Stop each container by name
	for _, containerName := range containerNames {
		log.Printf("Stopping container: %s", containerName)

		containerID, exists := containerMap[containerName]
		if !exists {
			errMsg := fmt.Sprintf("Container '%s' not found", containerName)
			log.Println(errMsg)
			errorMessages = append(errorMessages, errMsg)
			continue
		}

		// Set timeout for stopping the container (10 seconds)
		timeout := 10 * time.Second

		// Stop the container
		err := client.StopContainer(containerID, uint(timeout.Seconds()))
		if err != nil {
			errMsg := fmt.Sprintf("Failed to stop container '%s': %v", containerName, err)
			log.Println(errMsg)
			errorMessages = append(errorMessages, errMsg)
			// Continue to the next container even if one fails
			continue
		}

		log.Printf("Successfully stopped container '%s'", containerName)
	}

	if len(errorMessages) > 0 {
		return fmt.Errorf("encountered errors stopping one or more containers:\n%s",
			strings.Join(errorMessages, "\n"))
	}

	log.Println("All configured containers stopped (or attempted to stop) successfully.")
	return nil
}


