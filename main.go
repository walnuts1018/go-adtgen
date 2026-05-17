package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/walnuts1018/go-adtgen/internal/composer"
	"github.com/walnuts1018/go-adtgen/internal/emitter"
	"github.com/walnuts1018/go-adtgen/internal/loader"
	"github.com/walnuts1018/go-adtgen/internal/model"
	"github.com/walnuts1018/go-adtgen/internal/parser"
	"github.com/walnuts1018/go-adtgen/internal/resolver"
	"github.com/walnuts1018/go-adtgen/internal/writer"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	pattern := "."
	if len(args) > 0 {
		pattern = args[0]
	}

	pkg, err := loader.LoadGeneratePackage(loader.Config{Pattern: pattern})
	if err != nil {
		return err
	}

	declarations, err := parser.CollectDeclarations(pkg.Fset, pkg.SyntaxFiles())
	if err != nil {
		return err
	}

	resolved, err := resolver.ResolveDeclarations(pkg, declarations)
	if err != nil {
		return err
	}

	ordered, err := composer.OrderDeclarations(resolved)
	if err != nil {
		return err
	}

	generated := make([]model.GeneratedType, 0, len(ordered))
	for _, declaration := range ordered {
		generatedType, err := composer.BuildGeneratedType(declaration)
		if err != nil {
			return err
		}
		generated = append(generated, generatedType)
	}

	src, err := emitter.RenderForPackage(pkg.Package.PkgPath, pkg.Package.Name, generated)
	if err != nil {
		return err
	}

	output, err := outputPath(pkg)
	if err != nil {
		return err
	}

	return writer.WriteFile(output, src)
}

func outputPath(pkg *loader.Package) (string, error) {
	files := pkg.SyntaxFiles()
	if len(files) == 0 {
		return "", fmt.Errorf("generator: no syntax files loaded")
	}

	filename := pkg.Fset.Position(files[0].Package).Filename
	if filename == "" {
		return "", fmt.Errorf("generator: could not determine package directory")
	}

	return filepath.Join(filepath.Dir(filename), "zz_generated.product_types.go"), nil
}
