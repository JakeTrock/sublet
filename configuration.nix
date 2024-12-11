{ config, pkgs, ... }:

let
  nix-search = import (pkgs.fetchFromGitHub {
    owner = "diamondburned";
    repo = "nix-search";
    # You'll need to replace these with the latest values
    rev = "main";
    sha256 = ""; # Replace with the actual SHA
  });
in
{
  environment.systemPackages = [
    # existing packages...
    nix-search
  ];
} 