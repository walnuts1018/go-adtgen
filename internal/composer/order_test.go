package composer

import (
	"strings"
	"testing"

	"github.com/walnuts1018/go-product-type/internal/model"
)

func TestOrderDeclarationsPlacesDependenciesFirst(t *testing.T) {
	declarations := []model.ResolvedDeclaration{
		{
			Declaration: model.Declaration{
				Name:       "ABC",
				Expression: "AB[T, U] C[V]",
			},
			Inputs: []model.ResolvedType{
				{Expr: "AB[T, U]"},
				{Expr: "C[V]"},
			},
		},
		{
			Declaration: model.Declaration{
				Name:       "AB",
				Expression: "A[T] B[U]",
			},
			Inputs: []model.ResolvedType{
				{Expr: "A[T]"},
				{Expr: "B[U]"},
			},
		},
	}

	ordered, err := OrderDeclarations(declarations)
	if err != nil {
		t.Fatalf("OrderDeclarations() error = %v", err)
	}

	if len(ordered) != 2 {
		t.Fatalf("len(ordered) = %d, want 2", len(ordered))
	}
	if ordered[0].Declaration.Name != "AB" {
		t.Fatalf("ordered[0] = %q, want %q", ordered[0].Declaration.Name, "AB")
	}
	if ordered[1].Declaration.Name != "ABC" {
		t.Fatalf("ordered[1] = %q, want %q", ordered[1].Declaration.Name, "ABC")
	}
}

func TestOrderDeclarationsRejectsCycles(t *testing.T) {
	declarations := []model.ResolvedDeclaration{
		{
			Declaration: model.Declaration{Name: "AB"},
			Inputs: []model.ResolvedType{
				{Expr: "BC[T]"},
			},
		},
		{
			Declaration: model.Declaration{Name: "BC"},
			Inputs: []model.ResolvedType{
				{Expr: "AB[T]"},
			},
		},
	}

	_, err := OrderDeclarations(declarations)
	if err == nil {
		t.Fatal("OrderDeclarations() error = nil, want cycle error")
	}
	if !strings.Contains(err.Error(), "cycle") {
		t.Fatalf("OrderDeclarations() error = %q, want substring %q", err.Error(), "cycle")
	}
}

func TestOrderDeclarationsIgnoresQualifiedReferences(t *testing.T) {
	declarations := []model.ResolvedDeclaration{
		{
			Declaration: model.Declaration{Name: "AB"},
			Inputs: []model.ResolvedType{
				{Expr: "BC[T]"},
			},
		},
		{
			Declaration: model.Declaration{Name: "BC"},
			Inputs: []model.ResolvedType{
				{Expr: "pkg.AB[T]"},
			},
		},
	}

	ordered, err := OrderDeclarations(declarations)
	if err != nil {
		t.Fatalf("OrderDeclarations() error = %v", err)
	}

	if len(ordered) != 2 {
		t.Fatalf("len(ordered) = %d, want 2", len(ordered))
	}
	if ordered[0].Declaration.Name != "BC" {
		t.Fatalf("ordered[0] = %q, want %q", ordered[0].Declaration.Name, "BC")
	}
	if ordered[1].Declaration.Name != "AB" {
		t.Fatalf("ordered[1] = %q, want %q", ordered[1].Declaration.Name, "AB")
	}
}

func TestOrderDeclarationsSortsReadyDeclarationsDeterministically(t *testing.T) {
	declarations := []model.ResolvedDeclaration{
		{
			Declaration: model.Declaration{Name: "ZZ"},
			Inputs: []model.ResolvedType{
				{Expr: "AA[T]"},
				{Expr: "MM[T]"},
			},
		},
		{Declaration: model.Declaration{Name: "MM"}},
		{Declaration: model.Declaration{Name: "AA"}},
	}

	ordered, err := OrderDeclarations(declarations)
	if err != nil {
		t.Fatalf("OrderDeclarations() error = %v", err)
	}

	if len(ordered) != 3 {
		t.Fatalf("len(ordered) = %d, want 3", len(ordered))
	}
	if ordered[0].Declaration.Name != "AA" {
		t.Fatalf("ordered[0] = %q, want %q", ordered[0].Declaration.Name, "AA")
	}
	if ordered[1].Declaration.Name != "MM" {
		t.Fatalf("ordered[1] = %q, want %q", ordered[1].Declaration.Name, "MM")
	}
	if ordered[2].Declaration.Name != "ZZ" {
		t.Fatalf("ordered[2] = %q, want %q", ordered[2].Declaration.Name, "ZZ")
	}
}

func TestOrderDeclarationsRejectsDuplicateNames(t *testing.T) {
	declarations := []model.ResolvedDeclaration{
		{Declaration: model.Declaration{Name: "AB"}},
		{Declaration: model.Declaration{Name: "AB"}},
	}

	_, err := OrderDeclarations(declarations)
	if err == nil {
		t.Fatal("OrderDeclarations() error = nil, want duplicate name error")
	}
	if !strings.Contains(err.Error(), "duplicate") {
		t.Fatalf("OrderDeclarations() error = %q, want substring %q", err.Error(), "duplicate")
	}
}

func TestOrderDeclarationsIgnoresDeclarationTypeParameterShadowing(t *testing.T) {
	declarations := []model.ResolvedDeclaration{
		{
			Declaration: model.Declaration{
				Name:           "Foo",
				TypeParameters: []string{"T any"},
			},
			Inputs: []model.ResolvedType{
				{Expr: "Bar[T]"},
			},
		},
		{
			Declaration: model.Declaration{Name: "Bar"},
		},
		{
			Declaration: model.Declaration{Name: "T"},
			Inputs: []model.ResolvedType{
				{Expr: "Foo[Bar]"},
			},
		},
	}

	ordered, err := OrderDeclarations(declarations)
	if err != nil {
		t.Fatalf("OrderDeclarations() error = %v", err)
	}

	if len(ordered) != 3 {
		t.Fatalf("len(ordered) = %d, want 3", len(ordered))
	}
	if ordered[0].Declaration.Name != "Bar" {
		t.Fatalf("ordered[0] = %q, want %q", ordered[0].Declaration.Name, "Bar")
	}
	if ordered[1].Declaration.Name != "Foo" {
		t.Fatalf("ordered[1] = %q, want %q", ordered[1].Declaration.Name, "Foo")
	}
	if ordered[2].Declaration.Name != "T" {
		t.Fatalf("ordered[2] = %q, want %q", ordered[2].Declaration.Name, "T")
	}
}
