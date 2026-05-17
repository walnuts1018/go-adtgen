package composer

import (
	"go/token"
	"go/types"
	"testing"

	"github.com/walnuts1018/go-adtgen/internal/model"
)

func TestBuildGeneratedTypeGenericCollision(t *testing.T) {
	aType1 := types.NewNamed(
		types.NewTypeName(token.NoPos, nil, "A", nil),
		types.NewStruct(nil, nil),
		nil,
	)
	aType2 := types.NewNamed(
		types.NewTypeName(token.NoPos, nil, "A", nil),
		types.NewStruct(nil, nil),
		nil,
	)

	decl := model.ResolvedDeclaration{
		Declaration: model.Declaration{Name: "BothA"},
		Inputs: []model.ResolvedType{
			{Expr: "A[int]", Type: aType1, Struct: aType1.Underlying().(*types.Struct)},
			{Expr: "A[string]", Type: aType2, Struct: aType2.Underlying().(*types.Struct)},
		},
	}

	generated, err := BuildGeneratedType(decl)
	if err != nil {
		t.Fatal(err)
	}
	if generated.Inputs[0].ParameterName == generated.Inputs[1].ParameterName {
		t.Fatalf("Duplicate parameter names: %s", generated.Inputs[0].ParameterName)
	}
}
