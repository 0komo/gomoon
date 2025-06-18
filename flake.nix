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

  outputs = { self, flakelight, ... }@inputs:
    flakelight ./. {
      inherit inputs;
      withOverlays = [
        inputs.gomod2nix.overlays.default
      ];
      
      devShell = {
        packages = pkgs: with pkgs; [
          go
          gopls
          gomod2nix
        ];
      };
    };
}
