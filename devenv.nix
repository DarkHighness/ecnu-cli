{ pkgs, ... }:

{
  packages = [ pkgs.git pkgs.goreleaser ];

  languages.go.enable = true;
}
