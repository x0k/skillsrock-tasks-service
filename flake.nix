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
      lib = nixpkgs.lib;
      pkgs = import nixpkgs { inherit system; };
      unstablePkgs = import nixpkgsUnstable { inherit system; };
      buildGoModule = pkgs.buildGoModule.override {
        go = unstablePkgs.go_1_24;
      };
      mockery = buildGoModule rec {
        pname = "mockery";
        version = "2.53.3";

        src = pkgs.fetchFromGitHub {
          owner = "vektra";
          repo = "mockery";
          rev = "v${version}";
          sha256 = "sha256-X0cHpv4o6pzgjg7+ULCuFkspeff95WFtJbVHqy4LxAg=";
        };

        preCheck = ''
          substituteInPlace ./pkg/generator_test.go --replace-fail 0.0.0-dev ${version}
          substituteInPlace ./pkg/logging/logging_test.go --replace-fail v0.0 v${lib.versions.majorMinor version}
        '';

        ldflags = [
          "-s"
          "-w"
          "-X"
          "github.com/vektra/mockery/v2/pkg/logging.SemVer=v${version}"
        ];

        # env.CGO_ENABLED = false;

        proxyVendor = true;
        vendorHash = "sha256-AQY4x2bLqMwHIjoKHzEm1hebR29gRs3LJN8i00Uup5o=";

        subPackages = [ "." ];
      };
    in
    {
      devShells.${system} = {
        default = pkgs.mkShell {
          buildInputs = [
            mk.packages.${system}.default
            unstablePkgs.go_1_24
            mockery
            pkgs.go-migrate
            pkgs.golangci-lint
            pkgs.gotests
            pkgs.sqlc
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
