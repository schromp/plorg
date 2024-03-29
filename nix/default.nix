{
  pkgs ? (
    let
      sources = import ./nix/sources.nix;
    in
      import sources.nixpkgs {
        overlays = [
          (import "${sources.gomod2nix}/overlay.nix")
        ];
      }
  ),
  buildGoApplication,
  ...
}:
buildGoApplication {
  pname = "plorg";
  version = "0.1";
  pwd = ./.;
  src = ../.;
  modules = ./gomod2nix.toml;
}
