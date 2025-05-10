{
  description = "Container Refresh - A service to pull container images and stop containers";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = {
    self,
    nixpkgs,
    flake-utils,
    ...
  }:
    flake-utils.lib.eachDefaultSystem (
      system: let
        pkgs = nixpkgs.legacyPackages.${system};

        container-refresh = pkgs.buildGoModule {
          pname = "container-refresh";
          version = "0.1.0";
          src = ./.;

          vendorHash = null; # No external dependencies
        };
      in {
        packages = {
          container-refresh = container-refresh;
          default = container-refresh;
        };
      }
    )
    // {
      nixosModules.default = {
        config,
        lib,
        pkgs,
        ...
      }: let
        cfg = config.services.container-refresh;
      in {
        options.services.container-refresh = {
          enable = lib.mkEnableOption "container-refresh service";

          package = lib.mkOption {
            type = lib.types.package;
            default = self.packages.${pkgs.system}.container-refresh;
            description = "The container-refresh package to use.";
          };

          tokenFile = lib.mkOption {
            type = lib.types.path;
            description = "Path to the token file containing only the token to be used by the application.";
          };

          executable = lib.mkOption {
            type = lib.types.str;
            default = "docker";
            description = "The exact path to the container runtime to use.";
          };

          images = lib.mkOption {
            type = lib.types.listOf lib.types.str;
            default = [];
            description = "List of container image tags to pull.";
            example = lib.literalExpression ''              [
                            "docker.io/library/nginx:latest",
                            "ghcr.io/home-assistant/home-assistant:latest"
                          ]'';
          };

          containerNames = lib.mkOption {
            type = lib.types.listOf lib.types.str;
            default = [];
            description = "List of container names to stop after pulling images.";
            example = lib.literalExpression ''              [
                            "container1",
                            "container2"
                          ]'';
          };

          user = lib.mkOption {
            type = lib.types.str;
            default = "container-refresh";
            description = "User to run the service as.";
          };

          group = lib.mkOption {
            type = lib.types.str;
            default = "container-refresh";
            description = "Group to run the service as.";
          };

          port = lib.mkOption {
            type = lib.types.str;
            default = "8080";
            description = "Port to listen on.";
          };

          containerGroup = lib.mkOption {
            type = lib.types.enum ["docker" "podman"];
            default = "docker";
            description = "Container runtime group to add the service user to.";
          };
        };

        config = lib.mkIf cfg.enable {
          users.users.${cfg.user} = {
            isSystemUser = true;
            group = cfg.group;
            description = "container-refresh service user";
            home = "/var/lib/container-refresh";
            createHome = true;
            extraGroups = [ cfg.containerGroup ];
          };

          users.groups.${cfg.group} = {
            name = cfg.group;
          };

          systemd.services.container-refresh = {
            description = "Container Refresh Service";
            wantedBy = ["multi-user.target"];
            after = ["network.target"];
            unitConfig = {
              RequiresMountsFor = "/var/run/docker.sock";
            };

            serviceConfig = {
              ExecStart = "${cfg.package}/bin/container-refresh";
              Restart = "on-failure";
              User = cfg.user;
              Group = cfg.group;

              # Security hardening
              CapabilityBoundingSet = ["CAP_NET_BIND_SERVICE"];
              AmbientCapabilities = ["CAP_NET_BIND_SERVICE"];
              NoNewPrivileges = true;
              ProtectSystem = "strict";
              ProtectHome = true;
              PrivateTmp = true;
              PrivateDevices = true;

              # Allow access to container socket
              SupplementaryGroups = ["${cfg.containerGroup}"];
              ReadWritePaths = ["/var/run/docker.sock"];

              # Environment setup
              Environment =
                [
                  "TOKEN_FILE=${cfg.tokenFile}"
                  "PORT=${cfg.port}"
                  "CONTAINER_EXECUTABLE=${cfg.executable}"
                ]
                ++ lib.optionals (cfg.images != []) [
                  "CONTAINERS='${builtins.toJSON cfg.images}'"
                ]
                ++ lib.optionals (cfg.containerNames != []) [
                  "CONTAINER_NAMES='${builtins.toJSON cfg.containerNames}'"
                ];
            };
          };
        };
      };
    };
}
