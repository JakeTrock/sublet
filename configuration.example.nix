{ config, pkgs, ... }:

let
  sublet = import (pkgs.fetchFromGitHub {
    owner = "jaketrock";
    repo = "sublet";
    # You'll need to replace these with the latest values
    rev = "main";
    sha256 = ""; # Replace with the actual SHA
  });
in
{
  environment.systemPackages = [
    # existing packages...
    sublet
  ];

  imports = [
    # Your other imports...
    (builtins.getFlake "github:jaketrock/sublet").nixosModules.${pkgs.system}.systemd-subletd
  ];

  # Enable the service
  services.subletd.enable = true;
} 