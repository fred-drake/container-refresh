# Container Refresh Service

A Go application that listens for a webhook to update Docker containers.

## Configuration

The application is configured through environment variables:

- `TOKEN_FILE`: Path to a file containing the authentication token
- `CONTAINERS`: JSON array of container image tags to pull (e.g., `'["docker.io/library/alpine:latest"]'`)
- `CONTAINER_NAMES`: JSON array of container names to stop after pulling images (e.g., `'["my-container"]'`)
- `PORT`: The port for the HTTP server to listen on (defaults to `8080`)
- `CONTAINER_EXECUTABLE`: The container engine command to use (`docker` or `podman`, defaults to `docker`)

## API Endpoints

- `POST /update`
  - Body: `{"token": "YOUR_CONFIGURED_TOKEN"}`
  - Action: If the token is valid, pulls configured Docker containers and stops specified containers
  - Responses:
    - `200 OK`: Successfully processed the request
    - `401 Unauthorized`: Invalid token
    - `400 Bad Request`: Malformed JSON or missing token
    - `500 Internal Server Error`: If any error occurs during container operations

## Setup & Running

### Standard Go Installation

1. Ensure Go is installed
2. Build the application:
   ```bash
   go build -o container-refresh ./cmd/container-refresh
   ```
3. Set up environment variables:
   ```bash
   echo "my-secret-token" > token.txt
   export TOKEN_FILE="./token.txt"
   export CONTAINERS='["docker.io/library/nginx:latest"]'
   export CONTAINER_NAMES='["my-nginx"]'
   export PORT="8080"
   ```
4. Run the application:
   ```bash
   ./container-refresh
   ```

### Development

```bash
go mod tidy
go run ./cmd/container-refresh/main.go
```

## Nix Integration

### Building with Nix

To build the container-refresh binary using Nix:

```bash
nix build
```

This will create a `result` symlink pointing to the built package.

### Development Environment

This project uses [devenv](https://devenv.sh/) for development environment management. To set up your development environment:

```bash
devenv up
```

### NixOS Integration

To deploy container-refresh as a service on a NixOS system, add the following to your `configuration.nix`:

```nix
{ config, pkgs, ... }:

{
  imports = [
    # Import the flake's NixOS module
    (builtins.getFlake "path:to/container-refresh").nixosModules.default
  ];

  # Enable and configure the service
  services.container-refresh = {
    enable = true;
    
    # Required: Path to a file containing only the token
    tokenFile = "/path/to/token/file";
    
    # Required: List of container images to pull
    images = [
      "docker.io/library/nginx:latest",
      "ghcr.io/home-assistant/home-assistant:latest"
    ];
    
    # Optional: List of container names to stop after pulling
    containerNames = [
      "nginx",
      "home-assistant"
    ];
    
    # Optional: Change the user/group that runs the service
    # user = "container-refresh";
    # group = "container-refresh";
    
    # Optional: Change the listening port
    # port = "8080";
    
    # Optional: Choose container runtime (docker or podman)
    # executable = "docker";
  };
}
```

## Flake Outputs

The flake provides the following outputs:

- `packages.<system>.container-refresh`: The built Go binary
- `packages.<system>.default`: Same as above (default package)
- `nixosModules.default`: NixOS module for deploying as a service

## Security Considerations

The NixOS service is configured with security hardening:
- Runs as a dedicated user with minimal privileges
- Uses capability bounding to limit what the service can do
- Implements various security features (ProtectSystem, PrivateTmp, etc.)
- Has access to the Docker socket for its core functionality

## Requirements

- Docker must be installed and the `docker` group must exist
- The user running the service needs access to the Docker socket
