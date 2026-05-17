//go:build goproducttype_generate

package base

//goproducttype:product A[T] B
type AB[T any] struct{}
