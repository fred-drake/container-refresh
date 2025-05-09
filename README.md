# Container Refresh Service

A Go application that listens for a webhook to update Docker containers and restart associated systemd services.

## Configuration

The application requires a configuration file located at `/etc/container-refresh.toml`.
See `sample.container-refresh.toml` for an example structure.

Key configuration options:

- `Token` (string, required): The secret token for authorizing requests.
- `Containers` (array of strings, required): A list of full Docker image names to pull (e.g., `"docker.io/library/alpine:latest"`).
- `SystemdServices` (array of strings, required): A list of systemd service names to restart.
- `ServerPort` (string, optional): The port for the HTTP server to listen on. Defaults to `8080`.
- `Executable` (string, optional): The container engine command to use (e.g., `"docker"`, `"podman"`). Defaults to `"docker"`.

## API Endpoints

- `POST /update`
  - Body: `{"token": "YOUR_CONFIGURED_TOKEN"}`
  - Action: If the token is valid, pulls configured Docker containers and restarts configured systemd services.
  - Responses:
    - `200 OK`: Successfully processed the request.
    - `401 Unauthorized`: Invalid token.
    - `400 Bad Request`: Malformed JSON or missing token.
    - `500 Internal Server Error`: If any error occurs during container pulling or service restarting.

## Running in NixOS

See [NIX.md](NIX.md) for instructions on how to run container-refresh in NixOS.

## Setup & Running

1.  Ensure Go is installed.
2.  Create `/etc/container-refresh.toml` with your configuration.
3.  Build the application:
    ```bash
    go build -o container-refresh ./cmd/container-refresh
    ```
4.  Run the application (likely as a systemd service itself for production):
    ```bash
    ./container-refresh
    ```

## Development

```bash
go mod tidy
go run ./cmd/container-refresh/main.go
```
