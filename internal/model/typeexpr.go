package model

import (
	"go/ast"
	"go/parser"
	"go/types"
)

type ResolvedType struct {
	Expr   string
	Type   types.Type
	Struct *types.Struct
}

type ResolvedDeclaration struct {
	Declaration Declaration
	Inputs      []ResolvedType
}

func (t ResolvedType) ReferencedDeclarationNames(localNames map[string]struct{}) ([]string, error) {
	return t.ReferencedDeclarationNamesExcluding(localNames, nil)
}

func (t ResolvedType) ReferencedDeclarationNamesExcluding(localNames map[string]struct{}, excludedNames map[string]struct{}) ([]string, error) {
	expr, err := parser.ParseExpr(t.Expr)
	if err != nil {
		return nil, err
	}

	seen := make(map[string]struct{})
	names := make([]string, 0)
	collectReferencedDeclarationNames(expr, localNames, excludedNames, seen, &names)
	return names, nil
}

func collectReferencedDeclarationNames(expr ast.Expr, localNames map[string]struct{}, excludedNames map[string]struct{}, seen map[string]struct{}, names *[]string) {
	ast.Inspect(expr, func(node ast.Node) bool {
		switch node := node.(type) {
		case nil:
			return false
		case *ast.Ident:
			appendReferencedDeclarationName(node.Name, localNames, excludedNames, seen, names)
			return false
		case *ast.SelectorExpr:
			return false
		default:
			return true
		}
	})
}

func appendReferencedDeclarationName(name string, localNames map[string]struct{}, excludedNames map[string]struct{}, seen map[string]struct{}, names *[]string) {
	if _, ok := localNames[name]; !ok {
		return
	}
	if _, ok := excludedNames[name]; ok {
		return
	}
	if _, ok := seen[name]; ok {
		return
	}
	seen[name] = struct{}{}
	*names = append(*names, name)
}
