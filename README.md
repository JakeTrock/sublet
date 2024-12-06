# sublet

## Installation

### Nix

You can install `subletd` like so:

```nix
let
  subletd = import (pkgs.fetchFromGitHub {
    owner = "diamondburned";
    repo = "subletd";
    rev = "<REV>";
    sha256 = "<SHA256>";
  });
in

{
  environment.systemPackages = [ subletd ];
}
```
