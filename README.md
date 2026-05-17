# go-adtgen

`go-adtgen` generates Algebraic Data Types (sum types and product types) from annotated empty struct declarations.

## Usage

1. Add a `//go:build adtgen_generate` file.
2. Add `//go:generate go run github.com/walnuts1018/go-adtgen` to a normal package file. (Assuming you want to run the root package or the appropriate command path)
3. Add `//adtgen:product ...` or `//adtgen:sum ...` above an empty struct declaration.
4. Run `go generate ./...`.
