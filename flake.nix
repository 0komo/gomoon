{
  description = "A very basic flake";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    flakelight = {
      url = "github:nix-community/flakelight";
      inputs.nixpkgs.follows = "nixpkgs";
    };
    gomod2nix = {
      url = "github:nix-community/gomod2nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
    treefmt-nix = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs =
    { self, flakelight, ... }@inputs:
    flakelight ./. {
      inherit inputs;
      withOverlays = [
        inputs.gomod2nix.overlays.default
      ];

      devShell = {
        packages =
          pkgs: with pkgs; [
            lua5_4
            pkg-config
            clang-tools
            clang # fix stdlib not found on clangd
            bear
            go
            gopls
            gomod2nix
          ];
        stdenv = { clangStdenv, ... }: clangStdenv;
      };

      perSystem =
        pkgs:
        let
          treefmt = (inputs.treefmt-nix.lib.evalModule pkgs ./nix/treefmt.nix);
        in
        {
          formatter = treefmt.config.build.wrapper;
          # disable for now
          # checks.formatting = treefmt.config.build.check;
        };

      # checks = pkgs:
      # let
      #   mkTestBin =
      #     version:

      #     ;
      # in
      # {

      # };
    };
}
