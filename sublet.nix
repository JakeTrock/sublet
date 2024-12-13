{ lib, pkgs, config, ... }:

with lib;

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
        ExecStart = "${pkgs.sublet-go}/bin/sublet ${config.services.subletd.userId}";
      };
    };

    # Package the Go program
    environment.systemPackages = [ pkgs.sublet-go ];
  };

  # Define the Go package
  packages.sublet-go = pkgs.stdenv.mkDerivation {
    pname = "sublet-go";
    version = "1.0";

    src = ./.;

    buildInputs = [ pkgs.go ];

    buildPhase = ''
      mkdir -p $out/bin
      go build -o $out/bin/sublet ${src}
    '';

    installPhase = ''
      install -m755 $out/bin/sublet $out/bin/
    '';
  };
}


# packages.sublet = pkgs.buildGoModule {
# 					pname = "sublet";
# 					version = self.rev or "unknown";
# 					src = self;

# 					vendorHash = "sha256-ms8G6uXrp32zzE4fYfEqlo9Exfp2DnwUsq+BCyasJRg=";
# 					goDeps = ./go.mod;

# 					meta = with pkgs.lib; {
# 						description = "an LLM interface for NixOS";
# 						homepage = "https://github.com/jaketrock/sublet";
# 						mainProgram = "sublet";
# 					};
# 				};
# }

# https://mtlynch.io/notes/nix-import-from-url/