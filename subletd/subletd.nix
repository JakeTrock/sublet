{ lib, pkgs, config, ... }:

with lib;
let
  sublet-go = pkgs.buildGoModule {
      pname = "sublet";
      version = "0.1.0";
      src = ./.;

      vendorHash = "sha256-ms8G6uXrp32zzE4fYfEqlo9Exfp2DnwUsq+BCyasJRg=";
      CGO_ENABLED = "0";
      ldflags = [
        "-s"
        "-w"
      ];

      meta = with pkgs.lib; {
        description = "an LLM interface for NixOS";
        homepage = "https://github.com/jaketrock/sublet";
        mainProgram = "sublet";
      };
    };
in
{
  options.services.subletd = {
    enable = mkOption {
      type = types.bool;
      default = false;
      description = "Whether to enable the subletd service.";
    };

    hostUrl = mkOption {
      type = types.str;
      default = "http://159.65.35.90:8020";
      description = "The URL of the host to connect to";
    };

    userId = mkOption {
      type = types.str;
      default = "sublet";
      description = "Name of the user to associate with the sublet daemon";
    };
  };

  config = mkIf config.services.subletd.enable {
    systemd.services.subletd = {
      description = "A service that runs the sublet daemon";
      wantedBy = [ "multi-user.target" ];
      serviceConfig = {
        ExecStart = "${sublet-go}/bin/subletd ${config.services.subletd.userId} ${config.services.subletd.hostUrl}";
        Type = "simple";
        Restart = "always";
        RestartSec = "10";
      };
    };

    environment.systemPackages = [ sublet-go ];

  };
}

# https://mtlynch.io/notes/nix-import-from-url/