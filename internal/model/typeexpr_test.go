package model

import (
	"reflect"
	"testing"
)

func TestResolvedTypeReferencedDeclarationNamesWalksNestedExpressions(t *testing.T) {
	resolvedType := ResolvedType{
		Expr: "func(AB[T], pkg.Ignore[U]) map[BC[V]]*DE[W]",
	}

	names, err := resolvedType.ReferencedDeclarationNamesExcluding(map[string]struct{}{
		"AB":     {},
		"BC":     {},
		"DE":     {},
		"Ignore": {},
	}, nil)
	if err != nil {
		t.Fatalf("ReferencedDeclarationNames() error = %v", err)
	}

	want := []string{"AB", "BC", "DE"}
	if !reflect.DeepEqual(names, want) {
		t.Fatalf("ReferencedDeclarationNames() = %v, want %v", names, want)
	}
}

func TestResolvedTypeReferencedDeclarationNamesExcludesDeclarationTypeParameters(t *testing.T) {
	resolvedType := ResolvedType{
		Expr: "Bar[T]",
	}

	names, err := resolvedType.ReferencedDeclarationNamesExcluding(map[string]struct{}{
		"Bar": {},
		"T":   {},
	}, map[string]struct{}{
		"T": {},
	})
	if err != nil {
		t.Fatalf("ReferencedDeclarationNames() error = %v", err)
	}

	want := []string{"Bar"}
	if !reflect.DeepEqual(names, want) {
		t.Fatalf("ReferencedDeclarationNames() = %v, want %v", names, want)
	}
}
