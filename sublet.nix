{ lib, pkgs, config, ... }:

with lib;

{
  options.services.subletd = {
    enable = mkOption {
      type = types.bool;
      default = false;
      description = "Whether to enable the subletd service.";
    };

    helloTo = mkOption {
      type = types.str;
      default = "World";
      description = "Name of the person to greet.";
    };
  };

  config = mkIf config.services.subletd.enable {
    systemd.services.subletd = {
      description = "A service that prints a custom hello message";
      wantedBy = [ "multi-user.target" ];
      serviceConfig = {
        ExecStart = "${pkgs.coreutils}/bin/echo Hello, ${config.services.subletd.helloTo}!";
      };
    };
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