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
  options.services.subletd = {
    enable = mkOption {
      type = types.bool;
      default = false;
      description = "Whether to enable the sublet coordinator service.";
    };

    dbUrl = mkOption {
      type = types.str;
      default = "postgresql://postgres:postgres@localhost:5432/sublet";
      description = "The URL of the db to connect to";
    };
  };

  config = mkIf config.services.subletd.enable {
    systemd.services.subletd = {
      description = "A service that runs the sublet coordinator";
      wantedBy = [ "multi-user.target" ];
      serviceConfig = {
        ExecStart = "${sublet-node}/bin/server.js ${config.services.subletd.dbUrl}";
        WorkingDirectory = "${sublet-node}/bin";
        Type = "simple";
        Restart = "always";
        RestartSec = "10";
      };
    };

    environment.systemPackages = [ sublet-node ];

  };
}

# https://mtlynch.io/notes/nix-import-from-url/