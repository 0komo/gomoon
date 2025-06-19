{ ... }:
{
  projectRootFile = "flake.nix";
  programs = {
    gofmt.enable = true;
    nixfmt.enable = true;
    shellcheck.enable = true;
  };
}
