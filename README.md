# README

## About

Nix ai deploy

procedure:
Paste magic `curl sublet.sh | bash` into server, never look at it again
You get: monitoring, reverse proxy for testing, git listener so whenever agent alters git repo with nix script it rebuilds
Agent control panel lets you generate nix script verbally so you never need to configure server, lets you see liveness of servers and reverse proxy into them

Everything is versioned, reversible, and testable

based on https://github.com/elitak/nixos-infect
