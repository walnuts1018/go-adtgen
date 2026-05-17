# go-product-type

`go-product-type` generates product struct types from annotated empty struct declarations.

## Usage

1. Add a `//go:build adtgen_generate` file.
2. Add `//go:generate go run github.com/walnuts1018/go-adtgen/cmd/goproducttype` to a normal package file.
3. Add `//adtgen:product ...` above an empty struct declaration.
4. Run `go generate ./...`.
