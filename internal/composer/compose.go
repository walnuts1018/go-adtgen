package composer

import (
	"fmt"
	"go/types"

	"github.com/walnuts1018/go-product-type/internal/model"
)

type FieldSpec struct {
	Name      string
	Type      types.Type
	Tag       string
	Anonymous bool
}

func ExtractFields(st *types.Struct) []FieldSpec {
	if st == nil {
		return nil
	}

	fields := make([]FieldSpec, 0, st.NumFields())
	for i := 0; i < st.NumFields(); i++ {
		field := st.Field(i)
		fields = append(fields, FieldSpec{
			Name:      field.Name(),
			Type:      field.Type(),
			Tag:       st.Tag(i),
			Anonymous: field.Embedded(),
		})
	}
	return fields
}

func ComposeFields(groups [][]FieldSpec) ([]FieldSpec, error) {
	var composed []FieldSpec

	for _, group := range groups {
		for _, field := range group {
			index := findComposedField(composed, field)
			if index < 0 {
				composed = append(composed, field)
				continue
			}

			existing := composed[index]
			if existing.Anonymous != field.Anonymous {
				return nil, fmt.Errorf("composer: conflicting field %s", field.Name)
			}
			if field.Anonymous {
				// Identical anonymous fields merge by effective name/type, and the
				// first encountered field metadata (for example tags) is preserved.
				if !types.Identical(existing.Type, field.Type) {
					return nil, fmt.Errorf("composer: conflicting field %s", field.Name)
				}
				continue
			}

			if !types.Identical(existing.Type, field.Type) {
				return nil, fmt.Errorf("composer: conflicting field %s", field.Name)
			}
			if existing.Tag != field.Tag {
				return nil, fmt.Errorf("composer: conflicting tag for field %s", field.Name)
			}
		}
	}

	return composed, nil
}

func BuildGeneratedType(declaration model.ResolvedDeclaration) (model.GeneratedType, error) {
	groups := make([][]FieldSpec, 0, len(declaration.Inputs))
	for _, input := range declaration.Inputs {
		groups = append(groups, ExtractFields(input.Struct))
	}

	fields, err := ComposeFields(groups)
	if err != nil {
		return model.GeneratedType{}, err
	}

	generated := model.GeneratedType{
		Name:           declaration.Declaration.Name,
		TypeParameters: append([]string(nil), declaration.Declaration.TypeParameters...),
		Fields:         make([]model.GeneratedField, 0, len(fields)),
	}
	for _, field := range fields {
		generated.Fields = append(generated.Fields, model.GeneratedField{
			Name:      field.Name,
			Type:      field.Type,
			Tag:       field.Tag,
			Anonymous: field.Anonymous,
		})
	}

	return generated, nil
}

func findComposedField(fields []FieldSpec, target FieldSpec) int {
	for i, field := range fields {
		if field.Name == target.Name {
			return i
		}
	}
	return -1
}
