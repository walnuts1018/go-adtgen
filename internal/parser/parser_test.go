package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"
)

func TestCollectDeclarationsFindsProductAnnotation(t *testing.T) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "sample.go", `package sample
//goproducttype:product A B
type AB struct{}
`, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	decls, err := CollectDeclarations(fset, []*ast.File{file})
	if err != nil {
		t.Fatal(err)
	}
	if len(decls) != 1 {
		t.Fatalf("got %d declarations, want 1", len(decls))
	}
	if decls[0].Name != "AB" {
		t.Fatalf("got %q, want AB", decls[0].Name)
	}
	if decls[0].Expression != "A B" {
		t.Fatalf("got %q, want %q", decls[0].Expression, "A B")
	}
}

func TestCollectDeclarationsFindsTypeSpecAnnotationInGroupedTypeDeclaration(t *testing.T) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "sample.go", `package sample
type (
	//goproducttype:product A B
	AB struct{}
)
`, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	decls, err := CollectDeclarations(fset, []*ast.File{file})
	if err != nil {
		t.Fatal(err)
	}
	if len(decls) != 1 {
		t.Fatalf("got %d declarations, want 1", len(decls))
	}
	if decls[0].Name != "AB" {
		t.Fatalf("got %q, want AB", decls[0].Name)
	}
	if decls[0].Expression != "A B" {
		t.Fatalf("got %q, want %q", decls[0].Expression, "A B")
	}
}

func TestCollectDeclarationsIgnoresUnannotatedDeclarationsInGroupedTypeDeclaration(t *testing.T) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "sample.go", `package sample
type (
	//goproducttype:product A B
	AB struct{}
	CD struct{}
)
`, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	decls, err := CollectDeclarations(fset, []*ast.File{file})
	if err != nil {
		t.Fatal(err)
	}
	if len(decls) != 1 {
		t.Fatalf("got %d declarations, want 1", len(decls))
	}
	if decls[0].Name != "AB" {
		t.Fatalf("got %q, want AB", decls[0].Name)
	}
}

func TestCollectDeclarationsRejectsAnnotatedTypeAliasDeclarations(t *testing.T) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "sample.go", `package sample
//goproducttype:product A B
type AB = struct{}
`, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	_, err = CollectDeclarations(fset, []*ast.File{file})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "sample.go:3:6") {
		t.Fatalf("expected position-qualified error, got %q", err)
	}
	if !strings.Contains(err.Error(), "annotated declaration AB must not be a type alias") {
		t.Fatalf("unexpected error: %q", err)
	}
}

func TestCollectDeclarationsIgnoresSimilarDirectivePrefixes(t *testing.T) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "sample.go", `package sample
//goproducttype:productivity A B
type AB struct{}
`, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	decls, err := CollectDeclarations(fset, []*ast.File{file})
	if err != nil {
		t.Fatal(err)
	}
	if len(decls) != 0 {
		t.Fatalf("got %d declarations, want 0", len(decls))
	}
}

func TestCollectDeclarationsParsesWhitespaceTolerantDirective(t *testing.T) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "sample.go", "package sample\n//goproducttype:product\t  A B\ntype AB struct{}\n", parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	decls, err := CollectDeclarations(fset, []*ast.File{file})
	if err != nil {
		t.Fatal(err)
	}
	if len(decls) != 1 {
		t.Fatalf("got %d declarations, want 1", len(decls))
	}
	if decls[0].Expression != "A B" {
		t.Fatalf("got %q, want %q", decls[0].Expression, "A B")
	}
}
