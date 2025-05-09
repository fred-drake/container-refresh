# Nix Integration for Container Refresh

This document explains how to use the Nix flake to build and deploy the container-refresh service.

## Building with Nix

To build the container-refresh binary using Nix:

```bash
nix build
```

This will create a `result` symlink pointing to the built package.

## Development Shell

To enter a development shell with all necessary Go tools:

```bash
nix develop
```

## NixOS Integration

To deploy container-refresh as a systemd service on a NixOS system, add the following to your `configuration.nix`:

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
    
    # Optional: Override the default configuration path
    # configFile = "/etc/container-refresh.toml";
    
    # Optional: Change the user/group that runs the service
    # user = "container-refresh";
    # group = "container-refresh";
    
    # Optional: Change the listening port
    # port = "8080";
  };
}
```

### Configuration

Before starting the service, make sure to edit the configuration file at `/etc/container-refresh.toml` with your specific settings:

1. Set a secure token
2. Configure the containers to pull
3. Set up systemd services to restart

A sample configuration is provided in `sample.container-refresh.toml`.

## Flake Outputs

The flake provides the following outputs:

- `packages.<system>.container-refresh`: The built Go binary
- `packages.<system>.default`: Same as above (default package)
- `devShells.<system>.default`: Development environment with Go tools
- `nixosModules.default`: NixOS module for deploying as a systemd service

## Security Considerations

The systemd service is configured with security hardening:
- Runs as a dedicated user with minimal privileges
- Uses capability bounding to limit what the service can do
- Implements various systemd security features (ProtectSystem, PrivateTmp, etc.)
- Has access to the Docker socket and systemd journal for its core functionality

## Requirements

- Docker must be installed and the `docker` group must exist
- The user running the service needs access to the Docker socket
