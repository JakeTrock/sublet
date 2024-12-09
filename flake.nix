{
	inputs = {
		nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
		flake-utils.url = "github:numtide/flake-utils";
		flake-compat.url = "https://flakehub.com/f/edolstra/flake-compat/1.tar.gz";
		nix-search.url = "github:diamondburned/nix-search";
	};

	outputs = { self, nixpkgs, flake-utils, nix-search, ... }: flake-utils.lib.eachDefaultSystem (system:
		let
			pkgs = nixpkgs.legacyPackages.${system};
		in
		{
			devShells.default = pkgs.mkShell {
				packages = with pkgs; [
					go_1_23
					gopls
					gotools
					# sqlc
					nix-search.packages.${system}.default
				];
			};

			packages.default = pkgs.buildGoModule {
				pname = "subletd";
				version = self.rev or "unknown";
				src = self;

				vendorHash = "sha256-ms8G6uXrp32zzE4fYfEqlo9Exfp2DnwUsq+BCyasJRg=";
				goDeps = ./go.mod;

				meta = with pkgs.lib; {
					description = "A better and channel-compatible `nix search` for NixOS using only stable Nix tools.";
					homepage = "https://github.com/jaketrock/sublet";
					mainProgram = "subletd";
				};
			};

			apps = rec {
				default = subletd;
				subletd = {
					type = "app";
					program = "${self.packages.${system}.default}/bin/subletd";
				};
			};

			nixosModules = {
				systemd-subletd = { config, pkgs, lib, ... }: {
					options.services.subletd = lib.mkEnableOption "sublet daemon";

					config = lib.mkIf config.services.subletd.enable {
					systemd.services.subletd = {
						description = "sublet daemon";
						after = [ "network.target" ];
						wantedBy = [ "multi-user.target" ];
						serviceConfig = {
							ExecStart = "${self.packages.${system}.default}/bin/subletd";
							Restart = "always";
						};
					};
					};
				};
			};
		}
	);
}