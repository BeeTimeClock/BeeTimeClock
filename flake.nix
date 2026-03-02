{
  description = "Development environment with Python and Node.js";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-25.11";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let pkgs = nixpkgs.legacyPackages.${system};
      in {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            nodejs_24
            gcc
            libtool
            nodejs_22
            yarn
            typescript-language-server
            typescript
          ];

          shellHook = ''
            echo "🚀 Development environment loaded!"
            echo "📦 Node.js $(node --version)"
          '';
        };
      });
}
