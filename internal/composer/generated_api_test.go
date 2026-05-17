package composer

import (
	"go/token"
	"go/types"
	"testing"

	"github.com/walnuts1018/go-product-type/internal/model"
)

func TestBuildGeneratedTypeCreatesConstructorAndSplitMetadata(t *testing.T) {
	aType := types.NewNamed(
		types.NewTypeName(token.NoPos, nil, "A", nil),
		types.NewStruct([]*types.Var{
			types.NewField(token.NoPos, nil, "ID", types.Typ[types.String], false),
		}, []string{""}),
		nil,
	)
	bType := types.NewNamed(
		types.NewTypeName(token.NoPos, nil, "B", nil),
		types.NewStruct([]*types.Var{
			types.NewField(token.NoPos, nil, "Name", types.Typ[types.String], false),
		}, []string{""}),
		nil,
	)

	decl := model.ResolvedDeclaration{
		Declaration: model.Declaration{
			Name:           "AB",
			TypeParameters: []string{"T any"},
		},
		Inputs: []model.ResolvedType{
			{Expr: "A[T]", Type: aType, Struct: aType.Underlying().(*types.Struct)},
			{Expr: "B", Type: bType, Struct: bType.Underlying().(*types.Struct)},
		},
	}

	generated, err := BuildGeneratedType(decl)
	if err != nil {
		t.Fatal(err)
	}
	if len(generated.Inputs) != 2 {
		t.Fatalf("got %d inputs, want 2", len(generated.Inputs))
	}
	if generated.Inputs[0].ParameterName != "a" {
		t.Fatalf("got %q, want %q", generated.Inputs[0].ParameterName, "a")
	}
	if generated.Inputs[0].MethodName != "ToA" {
		t.Fatalf("got %q, want %q", generated.Inputs[0].MethodName, "ToA")
	}
	if generated.Inputs[1].ParameterName != "b" {
		t.Fatalf("got %q, want %q", generated.Inputs[1].ParameterName, "b")
	}
	if generated.Inputs[1].MethodName != "ToB" {
		t.Fatalf("got %q, want %q", generated.Inputs[1].MethodName, "ToB")
	}
}

func TestBuildGeneratedTypeDisambiguatesConflictingSplitMethodNames(t *testing.T) {
	leftPkg := types.NewPackage("example.com/left", "left")
	rightPkg := types.NewPackage("example.com/right", "right")

	leftType := types.NewNamed(types.NewTypeName(token.NoPos, leftPkg, "A", nil), types.NewStruct(nil, nil), nil)
	rightType := types.NewNamed(types.NewTypeName(token.NoPos, rightPkg, "A", nil), types.NewStruct(nil, nil), nil)

	decl := model.ResolvedDeclaration{
		Declaration: model.Declaration{Name: "BothA"},
		Inputs: []model.ResolvedType{
			{Expr: "left.A", Type: leftType, Struct: leftType.Underlying().(*types.Struct)},
			{Expr: "right.A", Type: rightType, Struct: rightType.Underlying().(*types.Struct)},
		},
	}

	generated, err := BuildGeneratedType(decl)
	if err != nil {
		t.Fatal(err)
	}
	if len(generated.Inputs) != 2 {
		t.Fatalf("got %d inputs, want 2", len(generated.Inputs))
	}
	if generated.Inputs[0].MethodName == generated.Inputs[1].MethodName {
		t.Fatal("expected distinct method names")
	}
}
