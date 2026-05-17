package model

import (
	"go/types"
	"reflect"
	"testing"
)

func TestGeneratedTypeKeepsInputMetadata(t *testing.T) {
	inputType := types.NewSlice(types.Typ[types.String])

	generated := GeneratedType{
		Name:           "AB",
		TypeParameters: []string{"T any"},
		Inputs: []GeneratedInput{
			{
				Expression:    "Foo[T]",
				ParameterName: "foo",
				Type:          inputType,
				MethodName:    "ToFoo",
				FieldNames:    []string{"ID", "Name"},
			},
		},
	}

	if len(generated.Inputs) != 1 {
		t.Fatalf("len(generated.Inputs) = %d, want 1", len(generated.Inputs))
	}

	input := generated.Inputs[0]
	if input.Expression != "Foo[T]" {
		t.Fatalf("input.Expression = %q, want %q", input.Expression, "Foo[T]")
	}
	if input.ParameterName != "foo" {
		t.Fatalf("input.ParameterName = %q, want %q", input.ParameterName, "foo")
	}
	if !types.Identical(input.Type, inputType) {
		t.Fatalf("input.Type = %v, want %v", input.Type, inputType)
	}
	if input.MethodName != "ToFoo" {
		t.Fatalf("input.MethodName = %q, want %q", input.MethodName, "ToFoo")
	}
	if !reflect.DeepEqual(input.FieldNames, []string{"ID", "Name"}) {
		t.Fatalf("input.FieldNames = %v, want %v", input.FieldNames, []string{"ID", "Name"})
	}
}
