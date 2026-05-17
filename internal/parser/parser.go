package parser

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"strings"

	"github.com/walnuts1018/go-product-type/internal/model"
)

const productDirective = "goproducttype:product"

func CollectDeclarations(fset *token.FileSet, files []*ast.File) ([]model.Declaration, error) {
	var declarations []model.Declaration

	for _, file := range files {
		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.TYPE {
				continue
			}

			for _, spec := range genDecl.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				expression, ok := productExpressionForTypeSpec(genDecl, typeSpec)
				if !ok {
					continue
				}
				if typeSpec.Assign.IsValid() {
					pos := fset.Position(typeSpec.Pos())
					return nil, fmt.Errorf("%s: annotated declaration %s must not be a type alias", pos, typeSpec.Name.Name)
				}
				structType, ok := typeSpec.Type.(*ast.StructType)
				if !ok || structType.Fields == nil || len(structType.Fields.List) != 0 {
					pos := fset.Position(typeSpec.Pos())
					return nil, fmt.Errorf("%s: annotated declaration %s must be an empty struct", pos, typeSpec.Name.Name)
				}

				declarations = append(declarations, model.Declaration{
					Name:           typeSpec.Name.Name,
					Expression:     expression,
					TypeParameters: collectTypeParameters(fset, typeSpec.TypeParams),
					Position:       fset.Position(typeSpec.Pos()),
				})
			}
		}
	}

	return declarations, nil
}

func productExpressionForTypeSpec(genDecl *ast.GenDecl, typeSpec *ast.TypeSpec) (string, bool) {
	if expression, ok := findProductExpression(typeSpec.Doc); ok {
		return expression, true
	}
	if expression, ok := findProductExpression(typeSpec.Comment); ok {
		return expression, true
	}
	if len(genDecl.Specs) == 1 {
		return findProductExpression(genDecl.Doc)
	}
	return "", false
}

func findProductExpression(group *ast.CommentGroup) (string, bool) {
	if group == nil {
		return "", false
	}

	for _, comment := range group.List {
		text := strings.TrimSpace(strings.TrimPrefix(comment.Text, "//"))
		fields := strings.Fields(text)
		if len(fields) == 0 {
			continue
		}
		directive := fields[0]
		if directive != productDirective {
			continue
		}
		expression := strings.TrimSpace(strings.TrimPrefix(text, directive))
		return expression, true
	}

	return "", false
}

func collectTypeParameters(fset *token.FileSet, fieldList *ast.FieldList) []string {
	if fieldList == nil {
		return nil
	}

	params := make([]string, 0, len(fieldList.List))
	for _, field := range fieldList.List {
		constraint := renderNode(fset, field.Type)
		for _, name := range field.Names {
			param := name.Name
			if constraint != "" {
				param += " " + constraint
			}
			params = append(params, param)
		}
	}

	return params
}

func renderNode(fset *token.FileSet, node ast.Node) string {
	if node == nil {
		return ""
	}

	var buf bytes.Buffer
	if err := format.Node(&buf, fset, node); err != nil {
		return ""
	}

	return buf.String()
}
