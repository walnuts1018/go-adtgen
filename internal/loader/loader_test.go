package loader

import (
	"go/types"
	"testing"
)

func TestLoadGeneratePackageRejectsNoPackagePattern(t *testing.T) {
	_, err := LoadGeneratePackage(Config{})
	if err == nil {
		t.Fatal("expected error for empty package pattern")
	}
}

func TestLoadGeneratePackageLoadsGeneratePackage(t *testing.T) {
	pkg, err := LoadGeneratePackage(Config{
		Pattern: "./testdata/loadgeneratepackage",
	})
	if err != nil {
		t.Fatalf("LoadGeneratePackage returned error: %v", err)
	}
	if pkg.Fset == nil {
		t.Fatal("expected file set")
	}
	if pkg.Package == nil {
		t.Fatal("expected loaded package")
	}
	if pkg.Package.Name != "loadgeneratepackage" {
		t.Fatalf("got package name %q, want %q", pkg.Package.Name, "loadgeneratepackage")
	}
	if len(pkg.Package.Syntax) != 1 {
		t.Fatalf("got %d syntax files, want 1", len(pkg.Package.Syntax))
	}
	if pkg.Package.Types == nil {
		t.Fatal("expected type information")
	}
	if pkg.Package.TypesInfo == nil {
		t.Fatal("expected types info")
	}
	if _, ok := pkg.Package.Types.Scope().Lookup("GenerateOnly").Type().Underlying().(*types.Struct); !ok {
		t.Fatal("expected GenerateOnly to resolve as a struct type")
	}
}

func TestPackageSyntaxFilesReturnsLoadedSyntax(t *testing.T) {
	pkg, err := LoadGeneratePackage(Config{
		Pattern: "./testdata/loadgeneratepackage",
	})
	if err != nil {
		t.Fatalf("LoadGeneratePackage returned error: %v", err)
	}

	files := pkg.SyntaxFiles()
	if len(files) != 1 {
		t.Fatalf("got %d syntax files, want 1", len(files))
	}
	if files[0] != pkg.Package.Syntax[0] {
		t.Fatal("expected returned syntax files to reference loaded package syntax")
	}
}
