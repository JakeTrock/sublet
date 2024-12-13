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

      installPhase = ''
        mkdir -p $out/bin
        cp $GOPATH/bin/sublet $out/bin/
      '';

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
        ExecStart = "${sublet-go}/bin/sublet ${config.services.subletd.userId}";
        Type = "simple";
        Restart = "always";
        RestartSec = "10";
      };
    };

    environment.systemPackages = [ sublet-go ];

  };
}

# https://mtlynch.io/notes/nix-import-from-url/