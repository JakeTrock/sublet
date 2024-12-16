{ pkgs, ... }:

let
#  nix-search = import (pkgs.fetchFromGitHub {
#    owner = "diamondburned";
#    repo = "nix-search";
#    rev = "e616ac1c82a616fa6e6d8c94839c5052eb8c808d";
#    hash = "sha256-h9yYOjL9i/m0r5NbqMcLMFNnwSKsIgfUr5qk+47pOtc=";
#  });
   subletrepo = builtins.fetchGit {
    url = "https://github.com/jaketrock/sublet";
    rev = "f06997419dc8db2effba9498564e2a4146bbe07d";
   };
in {
  imports = [
    ./hardware-configuration.nix
    ./networking.nix # generated at runtime by nixos-infect
    "${subletrepo}/subletd/subletd.nix" # daemon
  ];

  nix.settings.experimental-features = [ "nix-command" "flakes" ];

  boot.tmp.cleanOnBoot = true;
  zramSwap.enable = true;
  networking.hostName = "sublet01";
  networking.domain = "";
  services.openssh.enable = true;
  users.users.root.openssh.authorizedKeys.keys = [''ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIEZLyj8CP0LempUpevR1iGxoMhcVA/3oMJawA6HEgc/4 kapn'' ];
  environment.systemPackages = [ 
    pkgs.git
    pkgs.go
    pkgs.nodejs
#    nix-search
#    pkgs.nix-search-cli
#    sublet
  ];

  services.subletd = {
    enable = true;
    userId = "3dd437645316da7f7d63218de1b7c70b";
  };

  system.stateVersion = "24.11";
}