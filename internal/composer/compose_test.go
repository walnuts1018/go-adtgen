package composer

import (
	"go/token"
	"go/types"
	"strings"
	"testing"
)

func TestComposeFieldsMergesIdenticalNamedFields(t *testing.T) {
	fields := [][]FieldSpec{
		{
			{Name: "ID", Type: types.Typ[types.String], Tag: `json:"id"`},
			{Name: "Name", Type: types.Typ[types.String], Tag: `json:"name"`},
		},
		{
			{Name: "ID", Type: types.Typ[types.String], Tag: `json:"id"`},
		},
	}

	got, err := ComposeFields(fields)
	if err != nil {
		t.Fatalf("ComposeFields() error = %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("len(ComposeFields()) = %d, want 2", len(got))
	}
	if got[0].Name != "ID" {
		t.Fatalf("got first field %q, want %q", got[0].Name, "ID")
	}
	if got[1].Name != "Name" {
		t.Fatalf("got second field %q, want %q", got[1].Name, "Name")
	}
}

func TestComposeFieldsRejectsConflictingNamedFieldTypes(t *testing.T) {
	fields := [][]FieldSpec{
		{
			{Name: "ID", Type: types.Typ[types.String]},
		},
		{
			{Name: "ID", Type: types.Typ[types.Int]},
		},
	}

	_, err := ComposeFields(fields)
	if err == nil {
		t.Fatal("ComposeFields() error = nil, want conflict error")
	}
	if !strings.Contains(err.Error(), "conflicting field ID") {
		t.Fatalf("ComposeFields() error = %q, want substring %q", err.Error(), "conflicting field ID")
	}
}

func TestComposeFieldsRejectsConflictingNamedFieldTags(t *testing.T) {
	fields := [][]FieldSpec{
		{
			{Name: "ID", Type: types.Typ[types.String], Tag: `json:"id"`},
		},
		{
			{Name: "ID", Type: types.Typ[types.String], Tag: `db:"id"`},
		},
	}

	_, err := ComposeFields(fields)
	if err == nil {
		t.Fatal("ComposeFields() error = nil, want tag conflict error")
	}
	if !strings.Contains(err.Error(), "conflicting tag for field ID") {
		t.Fatalf("ComposeFields() error = %q, want substring %q", err.Error(), "conflicting tag for field ID")
	}
}

func TestComposeFieldsMergesIdenticalAnonymousFields(t *testing.T) {
	embeddedType := types.NewNamed(
		types.NewTypeName(token.NoPos, nil, "Embedded", nil),
		types.NewStruct(nil, nil),
		nil,
	)

	fields := [][]FieldSpec{
		{
			{Name: "Embedded", Type: embeddedType, Anonymous: true},
			{Name: "Name", Type: types.Typ[types.String]},
		},
		{
			{Name: "Embedded", Type: embeddedType, Anonymous: true},
		},
	}

	got, err := ComposeFields(fields)
	if err != nil {
		t.Fatalf("ComposeFields() error = %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("len(ComposeFields()) = %d, want 2", len(got))
	}
	if !got[0].Anonymous {
		t.Fatal("got[0].Anonymous = false, want true")
	}
	if got[0].Name != "Embedded" {
		t.Fatalf("got anonymous field %q, want %q", got[0].Name, "Embedded")
	}
}

func TestComposeFieldsKeepsFirstAnonymousFieldMetadataWhenIdenticalFieldsMerge(t *testing.T) {
	embeddedType := types.NewNamed(
		types.NewTypeName(token.NoPos, nil, "Embedded", nil),
		types.NewStruct(nil, nil),
		nil,
	)

	fields := [][]FieldSpec{
		{
			{Name: "Embedded", Type: embeddedType, Tag: `json:"embedded"`, Anonymous: true},
		},
		{
			{Name: "Embedded", Type: embeddedType, Tag: `yaml:"embedded"`, Anonymous: true},
		},
	}

	got, err := ComposeFields(fields)
	if err != nil {
		t.Fatalf("ComposeFields() error = %v, want nil", err)
	}

	if len(got) != 1 {
		t.Fatalf("len(ComposeFields()) = %d, want 1", len(got))
	}
	if got[0].Tag != `json:"embedded"` {
		t.Fatalf("got merged anonymous field tag %q, want first field tag %q", got[0].Tag, `json:"embedded"`)
	}
}

func TestComposeFieldsRejectsConflictingAnonymousFieldTypesWithSameEffectiveName(t *testing.T) {
	embeddedType := types.NewNamed(
		types.NewTypeName(token.NoPos, nil, "Embedded", nil),
		types.NewStruct(nil, nil),
		nil,
	)

	fields := [][]FieldSpec{
		{
			{Name: "Embedded", Type: embeddedType, Anonymous: true},
		},
		{
			{Name: "Embedded", Type: types.NewPointer(embeddedType), Anonymous: true},
		},
	}

	_, err := ComposeFields(fields)
	if err == nil {
		t.Fatal("ComposeFields() error = nil, want conflict error")
	}
	if !strings.Contains(err.Error(), "conflicting field Embedded") {
		t.Fatalf("ComposeFields() error = %q, want substring %q", err.Error(), "conflicting field Embedded")
	}
}

func TestComposeFieldsRejectsConflictingAnonymousFieldTypesFromDifferentPackages(t *testing.T) {
	leftType := types.NewNamed(
		types.NewTypeName(token.NoPos, types.NewPackage("example.com/left", "left"), "Embedded", nil),
		types.NewStruct(nil, nil),
		nil,
	)
	rightType := types.NewNamed(
		types.NewTypeName(token.NoPos, types.NewPackage("example.com/right", "right"), "Embedded", nil),
		types.NewStruct(nil, nil),
		nil,
	)

	fields := [][]FieldSpec{
		{
			{Name: "Embedded", Type: leftType, Anonymous: true},
		},
		{
			{Name: "Embedded", Type: rightType, Anonymous: true},
		},
	}

	_, err := ComposeFields(fields)
	if err == nil {
		t.Fatal("ComposeFields() error = nil, want conflict error")
	}
	if !strings.Contains(err.Error(), "conflicting field Embedded") {
		t.Fatalf("ComposeFields() error = %q, want substring %q", err.Error(), "conflicting field Embedded")
	}
}

func TestComposeFieldsRejectsConflictBetweenNamedAndAnonymousFieldsWithSameEffectiveName(t *testing.T) {
	embeddedType := types.NewNamed(
		types.NewTypeName(token.NoPos, nil, "Embedded", nil),
		types.NewStruct(nil, nil),
		nil,
	)

	fields := [][]FieldSpec{
		{
			{Name: "Embedded", Type: embeddedType, Anonymous: true},
		},
		{
			{Name: "Embedded", Type: types.Typ[types.String]},
		},
	}

	_, err := ComposeFields(fields)
	if err == nil {
		t.Fatal("ComposeFields() error = nil, want conflict error")
	}
	if !strings.Contains(err.Error(), "conflicting field Embedded") {
		t.Fatalf("ComposeFields() error = %q, want substring %q", err.Error(), "conflicting field Embedded")
	}
}

func TestComposeFieldsKeepsDistinctAnonymousFieldNamesForIdenticalTypes(t *testing.T) {
	embeddedType := types.NewNamed(
		types.NewTypeName(token.NoPos, nil, "Embedded", nil),
		types.NewStruct(nil, nil),
		nil,
	)
	aliasType := types.NewAlias(
		types.NewTypeName(token.NoPos, nil, "Alias", nil),
		embeddedType,
	)

	fields := [][]FieldSpec{
		{
			{Name: "Embedded", Type: embeddedType, Anonymous: true},
		},
		{
			{Name: "Alias", Type: aliasType, Anonymous: true},
		},
	}

	got, err := ComposeFields(fields)
	if err != nil {
		t.Fatalf("ComposeFields() error = %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("len(ComposeFields()) = %d, want 2", len(got))
	}
	if got[0].Name != "Embedded" {
		t.Fatalf("got first anonymous field %q, want %q", got[0].Name, "Embedded")
	}
	if got[1].Name != "Alias" {
		t.Fatalf("got second anonymous field %q, want %q", got[1].Name, "Alias")
	}
}

func TestExtractFieldsPreservesAnonymousFields(t *testing.T) {
	embeddedType := types.NewNamed(
		types.NewTypeName(token.NoPos, nil, "Embedded", nil),
		types.NewStruct([]*types.Var{
			types.NewField(token.NoPos, nil, "ID", types.Typ[types.String], false),
		}, []string{`json:"id"`}),
		nil,
	)

	st := types.NewStruct([]*types.Var{
		types.NewField(token.NoPos, nil, "Embedded", embeddedType, true),
	}, []string{""})

	fields := ExtractFields(st)
	if len(fields) != 1 {
		t.Fatalf("len(ExtractFields()) = %d, want 1", len(fields))
	}
	if !fields[0].Anonymous {
		t.Fatal("fields[0].Anonymous = false, want true")
	}
	if fields[0].Name != "Embedded" {
		t.Fatalf("fields[0].Name = %q, want %q", fields[0].Name, "Embedded")
	}
}

func TestExtractFieldsUsesOnlyDirectStructFields(t *testing.T) {
	embeddedType := types.NewNamed(
		types.NewTypeName(token.NoPos, nil, "Embedded", nil),
		types.NewStruct([]*types.Var{
			types.NewField(token.NoPos, nil, "Promoted", types.Typ[types.String], false),
		}, []string{`json:"promoted"`}),
		nil,
	)

	st := types.NewStruct([]*types.Var{
		types.NewField(token.NoPos, nil, "Embedded", embeddedType, true),
		types.NewField(token.NoPos, nil, "Name", types.Typ[types.String], false),
	}, []string{"", `json:"name"`})

	fields := ExtractFields(st)
	if len(fields) != 2 {
		t.Fatalf("len(ExtractFields()) = %d, want 2", len(fields))
	}
	if fields[0].Name != "Embedded" {
		t.Fatalf("fields[0].Name = %q, want %q", fields[0].Name, "Embedded")
	}
	if fields[1].Name != "Name" {
		t.Fatalf("fields[1].Name = %q, want %q", fields[1].Name, "Name")
	}
}
