package base

//go:generate go run ../../../../../cmd/goproducttype

type A[T any] struct {
	ID T
}

type B struct {
	Name string
}
