package loader

import (
	"errors"
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/packages"
)

type Config struct {
	Pattern string
}

type Package struct {
	Fset    *token.FileSet
	Package *packages.Package
}

func (p *Package) SyntaxFiles() []*ast.File {
	if p == nil || p.Package == nil {
		return nil
	}
	return p.Package.Syntax
}

func LoadGeneratePackage(cfg Config) (*Package, error) {
	if cfg.Pattern == "" {
		return nil, errors.New("loader: package pattern is required")
	}

	fset := token.NewFileSet()
	pkgs, err := packages.Load(&packages.Config{
		Fset: fset,
		Mode: packages.NeedName |
			packages.NeedSyntax |
			packages.NeedTypes |
			packages.NeedTypesInfo,
		BuildFlags: []string{"-tags=goproducttype_generate"},
	}, cfg.Pattern)
	if err != nil {
		return nil, err
	}
	if len(pkgs) != 1 {
		return nil, errors.New("loader: exactly one package must match the pattern")
	}

	pkg := pkgs[0]
	if len(pkg.Errors) > 0 {
		return nil, errors.New(formatPackageErrors(pkg))
	}
	return &Package{
		Fset:    fset,
		Package: pkg,
	}, nil
}

func formatPackageErrors(pkg *packages.Package) string {
	var b strings.Builder
	b.WriteString("loader: package load failed")
	for _, pkgErr := range pkg.Errors {
		b.WriteString(": ")
		b.WriteString(pkgErr.Error())
	}
	return b.String()
}
