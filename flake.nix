{
  description = "Tools for reading and manipulating PangYa game files";
  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  inputs.flake-utils.url = "github:numtide/flake-utils";
  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
      ...
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        inherit (pkgs) lib stdenv;
        pkgs = nixpkgs.legacyPackages.${system};
        deps = lib.optionals stdenv.isDarwin [ pkgs.macfuse-stubs ];
        gitCommit = self.dirtyShortRev or self.rev or "";
      in
      {
        packages = rec {
          pangfiles = pkgs.buildGoModule {
            name = "pangfiles";
            src = self;
            buildInputs = deps;
            vendorHash = pkgs.lib.fileContents ./go.mod.sri;
            ldflags = [ "-X github.com/pangbox/pangfiles/version.GitCommit=${gitCommit}" ];
            meta = {
              mainProgram = "pang";
            };
          };
          default = pangfiles;
        };
        devShells.default = pkgs.mkShell {
          packages = [
            pkgs.git
            pkgs.gopls
            pkgs.gotools
            pkgs.go
            pkgs.gnumake
            pkgs.nixfmt-rfc-style
          ]
          ++ deps;
        };
      }
    );
}
