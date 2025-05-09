{
  description = "A Blog AggreGATOR written in go :)";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";
  };

  outputs = {nixpkgs, ...}: {
    devShells = {
      x86_64-linux = let
        system = "x86_64-linux";
        pkgs = import nixpkgs {inherit system;};
      in {
        default = pkgs.mkShell {
          buildInputs = [pkgs.go pkgs.zsh];

          shellHook = ''
            export GOPATH=$PWD/.gopath
            export GOBIN=$GOPATH/bin
            export PATH=$GOBIN:$PATH
            mkdir -p "$GOBIN"
            go mod tidy
            go install github.com/pressly/goose/v3/cmd/goose@latest
            exec zsh
          '';
        };
      };
    };
  };
}
