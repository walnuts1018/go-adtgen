package composer

import (
	"go/types"
	"testing"

	"github.com/walnuts1018/go-adtgen/internal/model"
)

func TestBuildGeneratedTypeAnonymousStruct(t *testing.T) {
	structType := types.NewStruct(nil, nil)
	decl := model.ResolvedDeclaration{
		Declaration: model.Declaration{Name: "AB"},
		Inputs: []model.ResolvedType{
			{Expr: "struct{}", Type: structType, Struct: structType},
		},
	}
	generated, err := BuildGeneratedType(decl)
	if err != nil {
		t.Fatal(err)
	}
	if generated.Inputs[0].ParameterName != "struct{}" {
		t.Fatalf("Unexpected parameter name: %s", generated.Inputs[0].ParameterName)
	}
}
