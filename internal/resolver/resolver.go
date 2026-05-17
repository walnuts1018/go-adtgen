package resolver

import (
	"fmt"
	"go/ast"
	goparser "go/parser"
	"go/token"
	"go/types"
	"sort"
	"strings"
	"unicode"

	"github.com/walnuts1018/go-adtgen/internal/loader"
	"github.com/walnuts1018/go-adtgen/internal/model"
)

func ResolveExpression(expr string, files map[string]string) (types.Type, error) {
	fset := token.NewFileSet()

	parsedFiles, pkgPath, err := parseFiles(fset, files)
	if err != nil {
		return nil, err
	}

	conf := types.Config{}
	pkg, err := conf.Check(pkgPath, fset, parsedFiles, nil)
	if err != nil {
		return nil, err
	}

	resolved, err := resolveTypeInPackage(fset, pkg, expr)
	if err != nil {
		return nil, err
	}

	return resolved.Type, nil
}

func ResolveDeclarations(pkg *loader.Package, declarations []model.Declaration) ([]model.ResolvedDeclaration, error) {
	if pkg == nil || pkg.Fset == nil || pkg.Package == nil || pkg.Package.Types == nil {
		return nil, fmt.Errorf("resolver: loaded package with type information is required")
	}

	resolved := make([]model.ResolvedDeclaration, 0, len(declarations))
	for _, declaration := range declarations {
		declarationPos, err := declarationEvalPos(pkg, declaration)
		if err != nil {
			return nil, err
		}

		inputExprs, err := splitExpression(declaration.Expression)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", declaration.Position, err)
		}

		inputs := make([]model.ResolvedType, 0, len(inputExprs))
		for _, inputExpr := range inputExprs {
			input, err := resolveTypeInPackageAtPos(pkg.Fset, pkg.Package.Types, declarationPos, inputExpr)
			if err != nil {
				return nil, fmt.Errorf("%s: resolve %q: %w", declaration.Position, inputExpr, err)
			}
			if declaration.Kind == model.DeclarationKindSum {
				if err := validateSumInputType(pkg.Package.Types, input); err != nil {
					return nil, fmt.Errorf("%s: resolve %q: %w", declaration.Position, inputExpr, err)
				}
			}
			inputs = append(inputs, input)
		}

		resolved = append(resolved, model.ResolvedDeclaration{
			Declaration: declaration,
			Inputs:      inputs,
		})
	}

	return resolved, nil
}

func validateSumInputType(pkg *types.Package, input model.ResolvedType) error {
	if _, ok := input.Type.(*types.Alias); ok {
		return fmt.Errorf("sum inputs must be defined types, not aliases")
	}

	named, ok := types.Unalias(input.Type).(*types.Named)
	if !ok || named.Obj() == nil || named.Obj().Pkg() == nil {
		return fmt.Errorf("sum inputs must be defined in the same package")
	}
	if pkg == nil || named.Obj().Pkg().Path() != pkg.Path() {
		return fmt.Errorf("sum inputs must be defined in the same package")
	}
	return nil
}

func resolveTypeInPackage(fset *token.FileSet, pkg *types.Package, expr string) (model.ResolvedType, error) {
	return resolveTypeInPackageAtPos(fset, pkg, token.NoPos, expr)
}

func resolveTypeInPackageAtPos(fset *token.FileSet, pkg *types.Package, pos token.Pos, expr string) (model.ResolvedType, error) {
	tv, err := types.Eval(fset, pkg, pos, expr)
	if err != nil {
		return model.ResolvedType{}, err
	}
	if !tv.IsType() {
		return model.ResolvedType{}, fmt.Errorf("resolver: %q is not a type expression", expr)
	}
	structType, ok := structTypeOf(tv.Type)
	if !ok {
		return model.ResolvedType{}, fmt.Errorf("resolver: %q is not a struct type", expr)
	}

	return model.ResolvedType{
		Expr:   expr,
		Type:   tv.Type,
		Struct: structType,
	}, nil
}

func parseFiles(fset *token.FileSet, files map[string]string) ([]*ast.File, string, error) {
	names := make([]string, 0, len(files))
	for name := range files {
		names = append(names, name)
	}
	sort.Strings(names)

	parsedFiles := make([]*ast.File, 0, len(names))
	var pkgName string
	for _, name := range names {
		file, err := goparser.ParseFile(fset, name, files[name], goparser.ParseComments)
		if err != nil {
			return nil, "", err
		}
		if pkgName == "" {
			pkgName = file.Name.Name
		}
		parsedFiles = append(parsedFiles, file)
	}

	if pkgName == "" {
		return nil, "", fmt.Errorf("resolver: no files provided")
	}

	return parsedFiles, pkgName, nil
}

func structTypeOf(typ types.Type) (*types.Struct, bool) {
	if typ == nil {
		return nil, false
	}

	underlying := types.Unalias(typ).Underlying()
	structType, ok := underlying.(*types.Struct)
	return structType, ok
}

func declarationEvalPos(pkg *loader.Package, declaration model.Declaration) (token.Pos, error) {
	for _, file := range pkg.SyntaxFiles() {
		filename := pkg.Fset.Position(file.Pos()).Filename
		if filename != declaration.Position.Filename {
			continue
		}

		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.TYPE {
				continue
			}
			for _, spec := range genDecl.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok || typeSpec.Name.Name != declaration.Name {
					continue
				}

				pos := pkg.Fset.Position(typeSpec.Pos())
				if pos.Line == declaration.Position.Line && pos.Column == declaration.Position.Column {
					return typeSpec.Type.Pos(), nil
				}
			}
		}
	}

	return token.NoPos, fmt.Errorf("%s: could not locate declaration %s in loaded syntax", declaration.Position, declaration.Name)
}

func splitExpression(expr string) ([]string, error) {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return nil, fmt.Errorf("resolver: empty product expression")
	}

	parts := make([]string, 0, strings.Count(expr, " ")+1)
	start := -1
	depth := 0

	for i, r := range expr {
		if unicode.IsSpace(r) && depth == 0 {
			if start >= 0 {
				parts = append(parts, expr[start:i])
				start = -1
			}
			continue
		}

		if start < 0 {
			start = i
		}

		switch r {
		case '[':
			depth++
		case ']':
			depth--
			if depth < 0 {
				return nil, fmt.Errorf("resolver: malformed product expression %q", expr)
			}
		}
	}

	if depth != 0 {
		return nil, fmt.Errorf("resolver: malformed product expression %q", expr)
	}
	if start >= 0 {
		parts = append(parts, expr[start:])
	}
	if len(parts) == 0 {
		return nil, fmt.Errorf("resolver: empty product expression")
	}

	return parts, nil
}
