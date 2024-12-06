{
	inputs = {
		nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
		flake-utils.url = "github:numtide/flake-utils";
		flake-compat.url = "https://flakehub.com/f/edolstra/flake-compat/1.tar.gz";
	};

	outputs = { self, nixpkgs, flake-utils, ... }: flake-utils.lib.eachDefaultSystem (system:
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
		}
	);
}
