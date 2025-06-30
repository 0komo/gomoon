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
  };

  outputs =
    { self, flakelight, ... }@inputs:
    flakelight ./. (
      {
        lib,
        outputs,
        ...
      }:
      {
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
              treefmt
              gofumpt
              goimports-reviser
              nixfmt-rfc-style
              go
              gopls
              go-tools
              gomod2nix
            ];
          stdenv = { clangStdenv, ... }: clangStdenv;
        };

        checks = pkgs: {
          lint = pkgs.stdenvNoCC.mkDerivation {
            name = "lint";
            src = ./.;
            nativeBuildInputs = with pkgs; [
              go
              go-tools
            ];
            dontBuild = true;
            doCheck = true;
            checkPhase = ''
              export GOCACHE=$PWD/go-build
              export GOMODCACHE=$PWD/go/pkg/mod
              for tag in lua5{1..4}; do
                go vet -tags $tag
                staticcheck -tags $tag -f stylish
              done
            '';
          };
        };

        formatter = pkgs: pkgs.treefmt;
      }
    );
}
