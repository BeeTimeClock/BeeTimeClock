{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs = { self, nixpkgs }:
    let
      system = "x86_64-linux";
      pkgs = nixpkgs.legacyPackages.${system};
      fhs = pkgs.buildFHSEnv {
        name = "fhs-shell";
        targetPkgs = pkgs: [pkgs.gcc pkgs.libtool pkgs.nodejs_22 pkgs.yarn] ;
      };
    in
      {
        devShells.${system}.default = fhs.env;
      };
}
