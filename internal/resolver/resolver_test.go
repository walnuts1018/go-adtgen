package resolver

import (
	"go/types"
	"path/filepath"
	"strings"
	"testing"

	"github.com/walnuts1018/go-adtgen/internal/loader"
	"github.com/walnuts1018/go-adtgen/internal/model"
	"github.com/walnuts1018/go-adtgen/internal/parser"
)

const sampleNamedStructSource = `package sample
type A struct {
	Name string
}
`

const sampleFileName = "sample.go"

func TestResolveExpressionResolvesNamedStructType(t *testing.T) {
	files := map[string]string{
		sampleFileName: sampleNamedStructSource,
	}

	var (
		got types.Type
		err error
	)
	got, err = ResolveExpression("A", files)
	if err != nil {
		t.Fatalf("ResolveExpression returned error: %v", err)
	}
	if got == nil {
		t.Fatal("expected resolved type")
	}
	if _, ok := got.Underlying().(*types.Struct); !ok {
		t.Fatalf("got %T, want struct underlying type", got.Underlying())
	}
}

func TestResolveExpressionResolvesAliasToStructType(t *testing.T) {
	files := map[string]string{
		sampleFileName: `package sample
type Base struct {
	Name string
}

type Alias = Base
`,
	}

	got, err := ResolveExpression("Alias", files)
	if err != nil {
		t.Fatalf("ResolveExpression returned error: %v", err)
	}
	if got == nil {
		t.Fatal("expected resolved type")
	}
	if _, ok := got.Underlying().(*types.Struct); !ok {
		t.Fatalf("got %T, want struct underlying type", got.Underlying())
	}
}

func TestResolveExpressionResolvesTypeFromMultiFileSyntheticPackage(t *testing.T) {
	files := map[string]string{
		"a.go": sampleNamedStructSource,
		"b.go": `package sample
type B = A
`,
	}

	got, err := ResolveExpression("B", files)
	if err != nil {
		t.Fatalf("ResolveExpression returned error: %v", err)
	}
	if got == nil {
		t.Fatal("expected resolved type")
	}
	if _, ok := got.Underlying().(*types.Struct); !ok {
		t.Fatalf("got %T, want struct underlying type", got.Underlying())
	}
}

func TestResolveExpressionRejectsNonStructType(t *testing.T) {
	files := map[string]string{
		sampleFileName: `package sample
type A int
`,
	}

	_, err := ResolveExpression("A", files)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "not a struct type") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestResolveExpressionRejectsStructCompositeLiteral(t *testing.T) {
	files := map[string]string{
		sampleFileName: sampleNamedStructSource,
	}

	_, err := ResolveExpression("A{}", files)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "type expression") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestResolveExpressionRejectsStructTypedVariable(t *testing.T) {
	files := map[string]string{
		sampleFileName: `package sample
type A struct {
	Name string
}

var Value A
`,
	}

	_, err := ResolveExpression("Value", files)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "type expression") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestResolveDeclarationsResolvesLoadedPackageDeclarations(t *testing.T) {
	pkg, err := loader.LoadGeneratePackage(loader.Config{
		Pattern: filepath.Join("..", "testdata", "fixtures", "resolvepkg"),
	})
	if err != nil {
		t.Fatalf("LoadGeneratePackage returned error: %v", err)
	}

	decls, err := parser.CollectDeclarations(pkg.Fset, pkg.SyntaxFiles())
	if err != nil {
		t.Fatalf("CollectDeclarations returned error: %v", err)
	}

	resolved, err := ResolveDeclarations(pkg, decls)
	if err != nil {
		t.Fatalf("ResolveDeclarations returned error: %v", err)
	}

	assertResolvedDeclarationsFixture(t, resolved)
}

func assertResolvedDeclarationsFixture(t *testing.T, resolved []model.ResolvedDeclaration) {
	t.Helper()

	if len(resolved) != 5 {
		t.Fatalf("got %d resolved declarations, want 5", len(resolved))
	}

	resolvedByName := make(map[string]model.ResolvedDeclaration, len(resolved))
	for _, declaration := range resolved {
		resolvedByName[declaration.Declaration.Name] = declaration
	}

	assertProductDeclaration(t, resolvedByName, "CustomerAddress", "Customer Address", "Address")
	assertProductDeclaration(t, resolvedByName, "CustomerTime", "Customer timex.Time", "timex.Time")
	assertProductDeclaration(t, resolvedByName, "CustomerLocalTime", "Customer LocalTime", "LocalTime")
	assertGenericEnvelopeString(t, resolvedByName)
	assertGenericEnvelopeTypeParam(t, resolvedByName)
}

func assertProductDeclaration(t *testing.T, declarations map[string]model.ResolvedDeclaration, name, expression, secondExpr string) {
	t.Helper()

	declaration, ok := declarations[name]
	if !ok {
		t.Fatalf("expected %s declaration", name)
	}
	if declaration.Declaration.Expression != expression {
		t.Fatalf("got declaration expression %q, want %q", declaration.Declaration.Expression, expression)
	}
	if declaration.Declaration.Name != name {
		t.Fatalf("got declaration %q, want %q", declaration.Declaration.Name, name)
	}
	if len(declaration.Inputs) != 2 {
		t.Fatalf("got %d inputs, want 2", len(declaration.Inputs))
	}
	if declaration.Inputs[0].Expr != "Customer" {
		t.Fatalf("got first expr %q, want %q", declaration.Inputs[0].Expr, "Customer")
	}
	if declaration.Inputs[0].Type == nil || declaration.Inputs[0].Struct == nil {
		t.Fatal("expected first resolved input")
	}
	if declaration.Inputs[1].Expr != secondExpr {
		t.Fatalf("got second expr %q, want %q", declaration.Inputs[1].Expr, secondExpr)
	}
	if declaration.Inputs[1].Type == nil || declaration.Inputs[1].Struct == nil {
		t.Fatal("expected second resolved input")
	}
}

func assertGenericEnvelopeString(t *testing.T, declarations map[string]model.ResolvedDeclaration) {
	t.Helper()

	declaration := declarations["CustomerEnvelope"]
	assertProductDeclaration(t, declarations, "CustomerEnvelope", "Customer Envelope[string]", "Envelope[string]")
	if declaration.Inputs[1].Struct.NumFields() != 1 {
		t.Fatalf("got %d generic struct fields, want 1", declaration.Inputs[1].Struct.NumFields())
	}
	if declaration.Inputs[1].Struct.Field(0).Name() != "Value" {
		t.Fatalf("got generic field %q, want %q", declaration.Inputs[1].Struct.Field(0).Name(), "Value")
	}
	if got := types.TypeString(declaration.Inputs[1].Struct.Field(0).Type(), nil); got != "string" {
		t.Fatalf("got generic field type %q, want %q", got, "string")
	}
}

func assertGenericEnvelopeTypeParam(t *testing.T, declarations map[string]model.ResolvedDeclaration) {
	t.Helper()

	declaration := declarations["CustomerEnvelopeForTypeParam"]
	assertProductDeclaration(t, declarations, "CustomerEnvelopeForTypeParam", "Customer Envelope[T]", "Envelope[T]")
	if got := types.TypeString(declaration.Inputs[1].Struct.Field(0).Type(), nil); got != "T" {
		t.Fatalf("got type parameter field type %q, want %q", got, "T")
	}
}

func TestResolveDeclarationsResolvesSamePackageSumTypes(t *testing.T) {
	pkg, err := loader.LoadGeneratePackage(loader.Config{
		Pattern: filepath.Join("..", "testdata", "fixtures", "resolvepkgsumvalid"),
	})
	if err != nil {
		t.Fatalf("LoadGeneratePackage returned error: %v", err)
	}

	decls, err := parser.CollectDeclarations(pkg.Fset, pkg.SyntaxFiles())
	if err != nil {
		t.Fatalf("CollectDeclarations returned error: %v", err)
	}

	resolved, err := ResolveDeclarations(pkg, decls)
	if err != nil {
		t.Fatalf("ResolveDeclarations returned error: %v", err)
	}

	if len(resolved) != 1 {
		t.Fatalf("got %d resolved declarations, want 1", len(resolved))
	}
	if resolved[0].Declaration.Kind != model.DeclarationKindSum {
		t.Fatalf("got kind %q, want %q", resolved[0].Declaration.Kind, model.DeclarationKindSum)
	}
	if len(resolved[0].Inputs) != 2 {
		t.Fatalf("got %d inputs, want 2", len(resolved[0].Inputs))
	}
	if resolved[0].Inputs[0].Expr != "Hoge" {
		t.Fatalf("got first expr %q, want %q", resolved[0].Inputs[0].Expr, "Hoge")
	}
	if resolved[0].Inputs[1].Expr != "Fuga" {
		t.Fatalf("got second expr %q, want %q", resolved[0].Inputs[1].Expr, "Fuga")
	}
}

func TestResolveDeclarationsRejectsSumTypesFromExternalPackage(t *testing.T) {
	pkg, err := loader.LoadGeneratePackage(loader.Config{
		Pattern: filepath.Join("..", "testdata", "fixtures", "resolvepkgsumexternal"),
	})
	if err != nil {
		t.Fatalf("LoadGeneratePackage returned error: %v", err)
	}

	decls, err := parser.CollectDeclarations(pkg.Fset, pkg.SyntaxFiles())
	if err != nil {
		t.Fatalf("CollectDeclarations returned error: %v", err)
	}

	_, err = ResolveDeclarations(pkg, decls)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "sum inputs must be defined in the same package") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestResolveDeclarationsRejectsSumAliasInputs(t *testing.T) {
	pkg, err := loader.LoadGeneratePackage(loader.Config{
		Pattern: filepath.Join("..", "testdata", "fixtures", "resolvepkgsumalias"),
	})
	if err != nil {
		t.Fatalf("LoadGeneratePackage returned error: %v", err)
	}

	decls, err := parser.CollectDeclarations(pkg.Fset, pkg.SyntaxFiles())
	if err != nil {
		t.Fatalf("CollectDeclarations returned error: %v", err)
	}

	_, err = ResolveDeclarations(pkg, decls)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "sum inputs must be defined types, not aliases") {
		t.Fatalf("unexpected error: %v", err)
	}
}
