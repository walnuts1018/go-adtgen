# go-adtgen

`go-adtgen` generates Algebraic Data Types (sum types and product types) from annotated empty struct declarations.

## Usage

1. Add a `//go:build adtgen_generate` file.
2. Add the tool to your project:
   ```bash
   go get -tool github.com/walnuts1018/go-adtgen
   ```
3. Add `//go:generate go tool go-adtgen` to a normal package file.
4. Add `//adtgen:product ...` or `//adtgen:sum ...` above an empty struct declaration.
5. Run `go generate ./...`.
