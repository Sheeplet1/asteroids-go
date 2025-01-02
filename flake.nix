{
  description = "Nix-flake for Go development environment";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/release-24.11";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachSystem
      [
        "x86_64-linux"
        "aarch64-linux"
        "x86_64-darwin"
        "aarch64-darwin"
      ]
      (
        system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
        in
        {
          devShells.default = pkgs.mkShell {
            nativeBuildInputs = with pkgs; [
              wayland
              libxkbcommon
              xorg.libX11
              xorg.libXcursor
              xorg.libXi
              xorg.libXinerama
              xorg.libXfixes
              xorg.libXext
              xorg.libXrandr
              xorg.libXrender
              mesa
              libGL
            ];
            packages = with pkgs; [
              gopls
              gofumpt
              goimports-reviser
              gomodifytags
              golines
              gotests
            ];
          };
        }
      );
}
