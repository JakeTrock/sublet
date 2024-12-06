# sublet

## Installation

### Nix

You can install `subletd` like so:

```nix
let
  nix-search = import (pkgs.fetchFromGitHub {
    owner = "diamondburned";
    repo = "nix-search";
    rev = "<REV>";
    sha256 = "<SHA256>";
  });
in

{
  environment.systemPackages = [ nix-search ];
}
```
