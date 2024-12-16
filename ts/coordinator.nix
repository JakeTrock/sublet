{ lib, pkgs, config, ... }:

with lib;
let
  sublet-node = pkgs.stdenv.mkDerivation {
    pname = "sublet-coordinator";
    version = "0.1.0";
    src = ./.;

    buildInputs = with pkgs; [
      nodejs
      nodePackages.npm
    ];

    buildPhase = ''
      npm install
      npm run build
    '';

    installPhase = ''
      mkdir -p $out/bin
      cp -r dist/* $out/bin/
    '';

    meta = with pkgs.lib; {
      description = "A service that runs the sublet coordinator";
      homepage = "https://github.com/jaketrock/sublet";
    };
  };
in
{
  options.services.subletCoordinator = {
    enable = mkOption {
      type = types.bool;
      default = false;
      description = "Whether to enable the sublet coordinator service.";
    };

    user = mkOption {
      type = types.str;
      default = "sublet";
      description = "User account under which sublet runs.";
    };

    group = mkOption {
      type = types.str;
      default = "sublet";
      description = "Group under which sublet runs.";
    };

    dbName = mkOption {
      type = types.str;
      default = "sublet";
      description = "Name of the PostgreSQL database.";
    };

    dbUser = mkOption {
      type = types.str;
      default = "sublet";
      description = "PostgreSQL user for sublet.";
    };

    dbPassword = mkOption {
      type = types.str;
      default = "sublet";
      description = "PostgreSQL password for sublet user.";
    };
  };

  config = mkIf config.services.subletCoordinator.enable {
    # Create system user and group
    users.users.${config.services.subletCoordinator.user} = {
      isSystemUser = true;
      group = config.services.subletCoordinator.group;
      description = "Sublet coordinator service user";
    };

    users.groups.${config.services.subletCoordinator.group} = {};

    # Enable and configure PostgreSQL
    services.postgresql = {
      enable = true;
      ensureDatabases = [ config.services.subletCoordinator.dbName ];
      ensureUsers = [
        {
          name = config.services.subletCoordinator.dbUser;
          ensurePermissions = {
            "DATABASE ${config.services.subletCoordinator.dbName}" = "ALL PRIVILEGES";
          };
        }
      ];
    };

    # Configure systemd service
    systemd.services.subletCoordinator = {
      description = "Sublet coordinator service";
      wantedBy = [ "multi-user.target" ];
      after = [ "postgresql.service" ];
      requires = [ "postgresql.service" ];

      serviceConfig = {
        Type = "simple";
        User = config.services.subletCoordinator.user;
        Group = config.services.subletCoordinator.group;
        ExecStart = ''
          ${sublet-node}/bin/server.js \
          "postgresql://${config.services.subletCoordinator.dbUser}:${config.services.subletCoordinator.dbPassword}@localhost:5432/${config.services.subletCoordinator.dbName}"
        '';
        WorkingDirectory = "${sublet-node}/bin";
        Restart = "always";
        RestartSec = "10";
      };

      # Set up database password
      preStart = ''
        if ! ${pkgs.postgresql}/bin/psql -U ${config.services.subletCoordinator.dbUser} -d ${config.services.subletCoordinator.dbName} -c '\q' 2>/dev/null; then
          ${pkgs.postgresql}/bin/psql -U postgres -c "ALTER USER ${config.services.subletCoordinator.dbUser} WITH PASSWORD '${config.services.subletCoordinator.dbPassword}';"
        fi
      '';
    };

    environment.systemPackages = [ sublet-node ];
  };
}

# https://mtlynch.io/notes/nix-import-from-url/