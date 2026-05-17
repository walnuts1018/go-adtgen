package composer

import (
	"go/token"
	"go/types"
	"testing"

	"github.com/walnuts1018/go-product-type/internal/model"
)

func TestBuildGeneratedTypePreservesTypeParametersAndComposedFields(t *testing.T) {
	embeddedType := types.NewNamed(
		types.NewTypeName(token.NoPos, nil, "Embedded", nil),
		types.NewStruct(nil, nil),
		nil,
	)

	first := types.NewStruct(
		[]*types.Var{
			types.NewField(token.NoPos, nil, "ID", types.Typ[types.String], false),
			types.NewField(token.NoPos, nil, "Embedded", embeddedType, true),
		},
		[]string{`json:"id"`, `json:"embedded"`},
	)
	second := types.NewStruct(
		[]*types.Var{
			types.NewField(token.NoPos, nil, "ID", types.Typ[types.String], false),
			types.NewField(token.NoPos, nil, "Name", types.Typ[types.String], false),
		},
		[]string{`json:"id"`, `json:"name"`},
	)

	decl := model.ResolvedDeclaration{
		Declaration: model.Declaration{
			Name:           "AB",
			TypeParameters: []string{"T comparable", "U any"},
		},
		Inputs: []model.ResolvedType{
			{Struct: first},
			{Struct: second},
		},
	}

	generated, err := BuildGeneratedType(decl)
	if err != nil {
		t.Fatalf("BuildGeneratedType() error = %v", err)
	}

	if generated.Name != "AB" {
		t.Fatalf("generated.Name = %q, want %q", generated.Name, "AB")
	}
	if len(generated.TypeParameters) != 2 {
		t.Fatalf("len(generated.TypeParameters) = %d, want 2", len(generated.TypeParameters))
	}
	if generated.TypeParameters[0] != "T comparable" {
		t.Fatalf("generated.TypeParameters[0] = %q, want %q", generated.TypeParameters[0], "T comparable")
	}
	if generated.TypeParameters[1] != "U any" {
		t.Fatalf("generated.TypeParameters[1] = %q, want %q", generated.TypeParameters[1], "U any")
	}

	if len(generated.Fields) != 3 {
		t.Fatalf("len(generated.Fields) = %d, want 3", len(generated.Fields))
	}

	if generated.Fields[0].Name != "ID" {
		t.Fatalf("generated.Fields[0].Name = %q, want %q", generated.Fields[0].Name, "ID")
	}
	if generated.Fields[0].Tag != `json:"id"` {
		t.Fatalf("generated.Fields[0].Tag = %q, want %q", generated.Fields[0].Tag, `json:"id"`)
	}
	if generated.Fields[0].Anonymous {
		t.Fatal("generated.Fields[0].Anonymous = true, want false")
	}
	if !types.Identical(generated.Fields[0].Type, types.Typ[types.String]) {
		t.Fatalf("generated.Fields[0].Type = %v, want string", generated.Fields[0].Type)
	}

	if generated.Fields[1].Name != "Embedded" {
		t.Fatalf("generated.Fields[1].Name = %q, want %q", generated.Fields[1].Name, "Embedded")
	}
	if generated.Fields[1].Tag != `json:"embedded"` {
		t.Fatalf("generated.Fields[1].Tag = %q, want %q", generated.Fields[1].Tag, `json:"embedded"`)
	}
	if !generated.Fields[1].Anonymous {
		t.Fatal("generated.Fields[1].Anonymous = false, want true")
	}
	if !types.Identical(generated.Fields[1].Type, embeddedType) {
		t.Fatalf("generated.Fields[1].Type = %v, want %v", generated.Fields[1].Type, embeddedType)
	}

	if generated.Fields[2].Name != "Name" {
		t.Fatalf("generated.Fields[2].Name = %q, want %q", generated.Fields[2].Name, "Name")
	}
	if generated.Fields[2].Tag != `json:"name"` {
		t.Fatalf("generated.Fields[2].Tag = %q, want %q", generated.Fields[2].Tag, `json:"name"`)
	}
	if generated.Fields[2].Anonymous {
		t.Fatal("generated.Fields[2].Anonymous = true, want false")
	}
	if !types.Identical(generated.Fields[2].Type, types.Typ[types.String]) {
		t.Fatalf("generated.Fields[2].Type = %v, want string", generated.Fields[2].Type)
	}
}
