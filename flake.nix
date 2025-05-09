{
  description = "Container Refresh - A service to pull container images and restart systemd services";

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

          vendorHash = "sha256-CVycV7wxo7nOHm7qjZKfJrIkNcIApUNzN1mSIIwQN0g="; # Use the vendored dependencies
        };
      in {
        packages = {
          container-refresh = container-refresh;
          default = container-refresh;
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            gopls
            gotools
            go-tools
          ];
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

          configFile = lib.mkOption {
            type = lib.types.path;
            default = "/etc/container-refresh.toml";
            description = "Path to the configuration file.";
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
        };

        config = lib.mkIf cfg.enable {
          users.users.${cfg.user} = {
            isSystemUser = true;
            group = cfg.group;
            description = "container-refresh service user";
            home = "/var/lib/container-refresh";
            createHome = true;
          };

          users.groups.${cfg.group} = {};

          systemd.services.container-refresh = {
            description = "Container Refresh Service";
            wantedBy = ["multi-user.target"];
            after = ["network.target"];

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

              # Allow access to Docker socket and systemd
              SupplementaryGroups = ["docker" "systemd-journal"];
              ReadWritePaths = ["/var/run/docker.sock"];

              # Environment setup
              Environment = [
                "CONFIG_PATH=${cfg.configFile}"
              ];
            };
          };

          # Create a sample configuration file if it doesn't exist
          system.activationScripts.container-refresh-config = ''
            if [ ! -f ${cfg.configFile} ]; then
              cp ${./sample.container-refresh.toml} ${cfg.configFile}
              chmod 600 ${cfg.configFile}
              chown ${cfg.user}:${cfg.group} ${cfg.configFile}
            fi
          '';
        };
      };
    };
}
