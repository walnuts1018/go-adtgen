package base

import "testing"

func TestGeneratedAPIs(t *testing.T) {
	ab := NewAB[int](A[int]{ID: 1}, B{Name: "x"})
	if ab.ID != 1 || ab.Name != "x" {
		t.Fatal("unexpected constructor output")
	}
	if got := ab.ToA(); got.ID != 1 {
		t.Fatal("unexpected ToA result")
	}
	if got := ab.ToB(); got.Name != "x" {
		t.Fatal("unexpected ToB result")
	}
}
