# sublet

## Installation

### Nix

You can install `subletd` like so:

```nix
let
  subletd = import (pkgs.fetchFromGitHub {
    owner = "jaketrock";
    repo = "sublet";
    rev = "<REV>";
    sha256 = "<SHA256>";
  });
in

{
  environment.systemPackages = [ subletd ];
}
```
