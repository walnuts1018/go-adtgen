# go-product-type

`go-product-type` generates product struct types from annotated empty struct declarations.

## Usage

1. Add a `//go:build goproducttype_generate` file.
2. Add `//go:generate go run github.com/walnuts1018/go-product-type/cmd/goproducttype` to a normal package file.
3. Add `//goproducttype:product ...` above an empty struct declaration.
4. Run `go generate ./...`.
