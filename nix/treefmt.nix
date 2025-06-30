{ ... }:
{
  projectRootFile = "flake.nix";
  programs = {
    gofumpt.enable = true;
    nixfmt.enable = true;
    shellcheck.enable = true;
  };
}
