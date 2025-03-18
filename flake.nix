{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-24.11";
    nixpkgsUnstable.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    mk.url = "github:x0k/mk";
  };
  outputs =
    {
      self,
      nixpkgs,
      nixpkgsUnstable,
      mk,
    }:
    let
      system = "x86_64-linux";
      pkgs = import nixpkgs { inherit system; };
      unstablePkgs = import nixpkgsUnstable { inherit system; };
    in
    {
      devShells.${system} = {
        default = pkgs.mkShell {
          buildInputs = [
            mk.packages.${system}.default
            unstablePkgs.go_1_24
            pkgs.go-migrate
            pkgs.go-mockery
            pkgs.golangci-lint
            pkgs.gotests
            pkgs.delve
            pkgs.postgresql_17
          ];
          shellHook = ''
            source <(COMPLETE=bash mk)
          '';
          # CGO_CFLAGS="-U_FORTIFY_SOURCE -Wno-error";
          # CGO_CPPFLAGS="-U_FORTIFY_SOURCE -Wno-error";
        };
      };
    };
}
