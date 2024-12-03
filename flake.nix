{
#   description = "Sublet Go Server";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-23.11";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        packages.default = pkgs.buildGoModule {
          pname = "sublet-server";
          version = "0.1.0";
          src = ./go;  # Directory containing main.go and other Go files

          vendorHash = null;  # Will be replaced with actual hash on first build
        };

        # Development shell with required dependencies
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            gopls
            gotools
          ];
        };
      }
    ) // {
      # NixOS module for the systemd service
      nixosModules.default = { config, lib, pkgs, ... }:
        let
          cfg = config.services.sublet-server;
        in
        {
          options.services.sublet-server = {
            enable = lib.mkEnableOption "Sublet Go server";
            port = lib.mkOption {
              type = lib.types.port;
              default = 8020;
              description = "Port to listen on";
            };
          };

          config = lib.mkIf cfg.enable {
            systemd.services.sublet-server = {
            #   description = "Sublet Go Server";
              wantedBy = [ "multi-user.target" ];
              after = [ "network.target" ];

              serviceConfig = {
                ExecStart = "${self.packages.${pkgs.system}.default}/bin/sublet-server";
                Restart = "always";
                RestartSec = "10";
                Type = "simple";
                User = "sublet";
                Group = "sublet";
                WorkingDirectory = "/var/lib/sublet";

                # Hardening options
                NoNewPrivileges = true;
                ProtectSystem = "strict";
                ProtectHome = true;
                PrivateTmp = true;
                PrivateDevices = true;
                ProtectClock = true;
                ProtectProc = "invisible";
                ProtectKernelTunables = true;
                ProtectKernelModules = true;
                ProtectKernelLogs = true;
                ProtectControlGroups = true;
                RestrictAddressFamilies = [ "AF_INET" "AF_INET6" ];
                RestrictNamespaces = true;
                LockPersonality = true;
                MemoryDenyWriteExecute = true;
                RestrictRealtime = true;
                RestrictSUIDSGID = true;
                RemoveIPC = true;
                PrivateUsers = true;
                SystemCallArchitectures = "native";
              };
            };

            # Create user and group
            users.users.sublet = {
              isSystemUser = true;
              group = "sublet";
              home = "/var/lib/sublet";
              createHome = true;
            };

            users.groups.sublet = {};
          };
        };
    };
} 