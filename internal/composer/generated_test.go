package composer

import (
	"go/token"
	"go/types"
	"testing"

	"github.com/walnuts1018/go-adtgen/internal/model"
)

const (
	generatedFieldEmbedded   = "Embedded"
	generatedFieldName       = "Name"
	generatedTagJSONEmbedded = `json:"embedded"`
	generatedTagJSONName     = `json:"name"`
	generatedTypeHogeOrFuga  = "HogeOrFuga"
	generatedExprHoge        = "Hoge"
	generatedExprFuga        = "Fuga"
)

func TestBuildGeneratedTypePreservesTypeParametersAndComposedFields(t *testing.T) {
	embeddedType := types.NewNamed(
		types.NewTypeName(token.NoPos, nil, generatedFieldEmbedded, nil),
		types.NewStruct(nil, nil),
		nil,
	)

	first := types.NewStruct(
		[]*types.Var{
			types.NewField(token.NoPos, nil, "ID", types.Typ[types.String], false),
			types.NewField(token.NoPos, nil, generatedFieldEmbedded, embeddedType, true),
		},
		[]string{fieldTagJSONID, generatedTagJSONEmbedded},
	)
	second := types.NewStruct(
		[]*types.Var{
			types.NewField(token.NoPos, nil, "ID", types.Typ[types.String], false),
			types.NewField(token.NoPos, nil, generatedFieldName, types.Typ[types.String], false),
		},
		[]string{`json:"id"`, generatedTagJSONName},
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

	if generated.Fields[1].Name != generatedFieldEmbedded {
		t.Fatalf("generated.Fields[1].Name = %q, want %q", generated.Fields[1].Name, generatedFieldEmbedded)
	}
	if generated.Fields[1].Tag != generatedTagJSONEmbedded {
		t.Fatalf("generated.Fields[1].Tag = %q, want %q", generated.Fields[1].Tag, generatedTagJSONEmbedded)
	}
	if !generated.Fields[1].Anonymous {
		t.Fatal("generated.Fields[1].Anonymous = false, want true")
	}
	if !types.Identical(generated.Fields[1].Type, embeddedType) {
		t.Fatalf("generated.Fields[1].Type = %v, want %v", generated.Fields[1].Type, embeddedType)
	}

	if generated.Fields[2].Name != generatedFieldName {
		t.Fatalf("generated.Fields[2].Name = %q, want %q", generated.Fields[2].Name, generatedFieldName)
	}
	if generated.Fields[2].Tag != generatedTagJSONName {
		t.Fatalf("generated.Fields[2].Tag = %q, want %q", generated.Fields[2].Tag, generatedTagJSONName)
	}
	if generated.Fields[2].Anonymous {
		t.Fatal("generated.Fields[2].Anonymous = true, want false")
	}
	if !types.Identical(generated.Fields[2].Type, types.Typ[types.String]) {
		t.Fatalf("generated.Fields[2].Type = %v, want string", generated.Fields[2].Type)
	}
}

func TestBuildGeneratedTypeBuildsSumMetadataFromEmbeddedCommonFields(t *testing.T) {
	pkg := types.NewPackage("example.com/sample", "sample")
	commonType := types.NewNamed(
		types.NewTypeName(token.NoPos, pkg, "Common", nil),
		types.NewStruct([]*types.Var{
			types.NewField(token.NoPos, pkg, "ID", types.Typ[types.String], false),
		}, nil),
		nil,
	)
	hogeType := types.NewNamed(
		types.NewTypeName(token.NoPos, pkg, "Hoge", nil),
		types.NewStruct([]*types.Var{
			types.NewField(token.NoPos, pkg, "Common", commonType, true),
			types.NewField(token.NoPos, pkg, "Name", types.Typ[types.String], false),
		}, nil),
		nil,
	)
	fugaType := types.NewNamed(
		types.NewTypeName(token.NoPos, pkg, "Fuga", nil),
		types.NewStruct([]*types.Var{
			types.NewField(token.NoPos, pkg, "Common", commonType, true),
			types.NewField(token.NoPos, pkg, "Age", types.Typ[types.Int], false),
		}, nil),
		nil,
	)

	decl := model.ResolvedDeclaration{
		Declaration: model.Declaration{
			Kind: model.DeclarationKindSum,
			Name: generatedTypeHogeOrFuga,
		},
		InterfaceMethods: []model.ResolvedInterfaceMethod{
			{
				Name:      "String",
				Signature: types.NewSignatureType(nil, nil, nil, nil, types.NewTuple(types.NewVar(token.NoPos, nil, "", types.Typ[types.String])), false),
			},
		},
		Inputs: []model.ResolvedType{
			{Expr: generatedExprHoge, Type: hogeType, Struct: hogeType.Underlying().(*types.Struct)},
			{Expr: generatedExprFuga, Type: fugaType, Struct: fugaType.Underlying().(*types.Struct)},
		},
	}

	generated, err := BuildGeneratedType(decl)
	if err != nil {
		t.Fatalf("BuildGeneratedType() error = %v", err)
	}
	if generated.Kind != model.DeclarationKindSum {
		t.Fatalf("generated.Kind = %q, want %q", generated.Kind, model.DeclarationKindSum)
	}
	if generated.Sum == nil {
		t.Fatal("generated.Sum = nil")
	}
	if len(generated.Sum.InterfaceMethods) != 1 {
		t.Fatalf("len(generated.Sum.InterfaceMethods) = %d, want 1", len(generated.Sum.InterfaceMethods))
	}
	if generated.Sum.InterfaceMethods[0].Name != "String" {
		t.Fatalf("generated.Sum.InterfaceMethods[0].Name = %q, want %q", generated.Sum.InterfaceMethods[0].Name, "String")
	}
	if len(generated.Sum.Variants) != 2 {
		t.Fatalf("len(generated.Sum.Variants) = %d, want 2", len(generated.Sum.Variants))
	}
	if generated.Sum.Variants[0].TypeName != generatedExprHoge {
		t.Fatalf("generated.Sum.Variants[0].TypeName = %q, want %q", generated.Sum.Variants[0].TypeName, generatedExprHoge)
	}
	if generated.Sum.Variants[1].TypeName != generatedExprFuga {
		t.Fatalf("generated.Sum.Variants[1].TypeName = %q, want %q", generated.Sum.Variants[1].TypeName, generatedExprFuga)
	}
	if len(generated.Sum.CommonFields) != 1 {
		t.Fatalf("len(generated.Sum.CommonFields) = %d, want 1", len(generated.Sum.CommonFields))
	}
	field := generated.Sum.CommonFields[0]
	if field.Name != "ID" {
		t.Fatalf("field.Name = %q, want %q", field.Name, "ID")
	}
	if field.GetterName != "GetID" {
		t.Fatalf("field.GetterName = %q, want %q", field.GetterName, "GetID")
	}
	if field.SetterName != "SetID" {
		t.Fatalf("field.SetterName = %q, want %q", field.SetterName, "SetID")
	}
	if !generated.Sum.GenerateSetters {
		t.Fatal("generated.Sum.GenerateSetters = false, want true")
	}
	if len(field.Paths) != 2 {
		t.Fatalf("len(field.Paths) = %d, want 2", len(field.Paths))
	}
	if got := field.Paths[0]; len(got) != 2 || got[0] != "Common" || got[1] != "ID" {
		t.Fatalf("field.Paths[0] = %v, want [Common ID]", got)
	}
	if got := field.Paths[1]; len(got) != 2 || got[0] != "Common" || got[1] != "ID" {
		t.Fatalf("field.Paths[1] = %v, want [Common ID]", got)
	}
}

func TestBuildGeneratedTypeDisablesSumSettersWithOption(t *testing.T) {
	pkg := types.NewPackage("example.com/sample", "sample")
	hogeType := types.NewNamed(
		types.NewTypeName(token.NoPos, pkg, "Hoge", nil),
		types.NewStruct([]*types.Var{
			types.NewField(token.NoPos, pkg, "ID", types.Typ[types.String], false),
		}, nil),
		nil,
	)
	fugaType := types.NewNamed(
		types.NewTypeName(token.NoPos, pkg, "Fuga", nil),
		types.NewStruct([]*types.Var{
			types.NewField(token.NoPos, pkg, "ID", types.Typ[types.String], false),
		}, nil),
		nil,
	)

	decl := model.ResolvedDeclaration{
		Declaration: model.Declaration{
			Kind: model.DeclarationKindSum,
			Name: "HogeOrFuga",
			Options: model.DeclarationOptions{
				NoSetter: true,
			},
		},
		Inputs: []model.ResolvedType{
			{Expr: "Hoge", Type: hogeType, Struct: hogeType.Underlying().(*types.Struct)},
			{Expr: "Fuga", Type: fugaType, Struct: fugaType.Underlying().(*types.Struct)},
		},
	}

	generated, err := BuildGeneratedType(decl)
	if err != nil {
		t.Fatalf("BuildGeneratedType() error = %v", err)
	}
	if generated.Sum == nil {
		t.Fatal("generated.Sum = nil")
	}
	if generated.Sum.GenerateSetters {
		t.Fatal("generated.Sum.GenerateSetters = true, want false")
	}
	if len(generated.Sum.CommonFields) != 1 {
		t.Fatalf("len(generated.Sum.CommonFields) = %d, want 1", len(generated.Sum.CommonFields))
	}
	if generated.Sum.CommonFields[0].SetterName != "SetID" {
		t.Fatalf("generated.Sum.CommonFields[0].SetterName = %q, want %q", generated.Sum.CommonFields[0].SetterName, "SetID")
	}
}

func TestBuildGeneratedTypeRejectsAmbiguousPromotedCommonFields(t *testing.T) {
	pkg := types.NewPackage("example.com/sample", "sample")
	leftType := types.NewNamed(
		types.NewTypeName(token.NoPos, pkg, "Left", nil),
		types.NewStruct([]*types.Var{
			types.NewField(token.NoPos, pkg, "ID", types.Typ[types.String], false),
		}, nil),
		nil,
	)
	rightType := types.NewNamed(
		types.NewTypeName(token.NoPos, pkg, "Right", nil),
		types.NewStruct([]*types.Var{
			types.NewField(token.NoPos, pkg, "ID", types.Typ[types.String], false),
		}, nil),
		nil,
	)
	hogeType := types.NewNamed(
		types.NewTypeName(token.NoPos, pkg, "Hoge", nil),
		types.NewStruct([]*types.Var{
			types.NewField(token.NoPos, pkg, "Left", leftType, true),
			types.NewField(token.NoPos, pkg, "Right", rightType, true),
		}, nil),
		nil,
	)
	fugaType := types.NewNamed(
		types.NewTypeName(token.NoPos, pkg, "Fuga", nil),
		types.NewStruct([]*types.Var{
			types.NewField(token.NoPos, pkg, "ID", types.Typ[types.String], false),
		}, nil),
		nil,
	)

	decl := model.ResolvedDeclaration{
		Declaration: model.Declaration{
			Kind: model.DeclarationKindSum,
			Name: "HogeOrFuga",
		},
		Inputs: []model.ResolvedType{
			{Expr: "Hoge", Type: hogeType, Struct: hogeType.Underlying().(*types.Struct)},
			{Expr: "Fuga", Type: fugaType, Struct: fugaType.Underlying().(*types.Struct)},
		},
	}

	_, err := BuildGeneratedType(decl)
	if err == nil {
		t.Fatal("expected error")
	}
	if got := err.Error(); got != "composer: ambiguous promoted field ID in Hoge" {
		t.Fatalf("BuildGeneratedType() error = %q, want %q", got, "composer: ambiguous promoted field ID in Hoge")
	}
}

func TestBuildGeneratedTypeSkipsConflictingCommonFieldTypes(t *testing.T) {
	pkg := types.NewPackage("example.com/sample", "sample")
	hogeType := types.NewNamed(
		types.NewTypeName(token.NoPos, pkg, "Hoge", nil),
		types.NewStruct([]*types.Var{
			types.NewField(token.NoPos, pkg, "ID", types.Typ[types.String], false),
		}, nil),
		nil,
	)
	fugaType := types.NewNamed(
		types.NewTypeName(token.NoPos, pkg, "Fuga", nil),
		types.NewStruct([]*types.Var{
			types.NewField(token.NoPos, pkg, "ID", types.Typ[types.Int], false),
		}, nil),
		nil,
	)

	decl := model.ResolvedDeclaration{
		Declaration: model.Declaration{
			Kind: model.DeclarationKindSum,
			Name: "HogeOrFuga",
		},
		Inputs: []model.ResolvedType{
			{Expr: "Hoge", Type: hogeType, Struct: hogeType.Underlying().(*types.Struct)},
			{Expr: "Fuga", Type: fugaType, Struct: fugaType.Underlying().(*types.Struct)},
		},
	}

	generated, err := BuildGeneratedType(decl)
	if err != nil {
		t.Fatalf("BuildGeneratedType() error = %v", err)
	}
	if generated.Sum == nil {
		t.Fatal("generated.Sum = nil")
	}
	if len(generated.Sum.CommonFields) != 0 {
		t.Fatalf("len(generated.Sum.CommonFields) = %d, want 0", len(generated.Sum.CommonFields))
	}
}
