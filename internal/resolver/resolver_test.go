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

func TestResolveExpressionResolvesNamedStructType(t *testing.T) {
	files := map[string]string{
		"sample.go": `package sample
type A struct {
	Name string
}
`,
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
		"sample.go": `package sample
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
		"a.go": `package sample
type A struct {
	Name string
}
`,
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
		"sample.go": `package sample
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
		"sample.go": `package sample
type A struct {
	Name string
}
`,
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
		"sample.go": `package sample
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

	if len(resolved) != 5 {
		t.Fatalf("got %d resolved declarations, want 5", len(resolved))
	}

	resolvedByName := make(map[string]int, len(resolved))
	for i, declaration := range resolved {
		resolvedByName[declaration.Declaration.Name] = i
	}

	samePackage, ok := resolvedByName["CustomerAddress"]
	if !ok {
		t.Fatal("expected CustomerAddress declaration")
	}
	genericIndex, ok := resolvedByName["CustomerEnvelope"]
	if !ok {
		t.Fatal("expected CustomerEnvelope declaration")
	}
	importedIndex, ok := resolvedByName["CustomerTime"]
	if !ok {
		t.Fatal("expected CustomerTime declaration")
	}
	aliasIndex, ok := resolvedByName["CustomerLocalTime"]
	if !ok {
		t.Fatal("expected CustomerLocalTime declaration")
	}
	typeParamIndex, ok := resolvedByName["CustomerEnvelopeForTypeParam"]
	if !ok {
		t.Fatal("expected CustomerEnvelopeForTypeParam declaration")
	}

	samePackageDecl := resolved[samePackage]
	if samePackageDecl.Declaration.Expression != "Customer Address" {
		t.Fatalf("got declaration expression %q, want %q", samePackageDecl.Declaration.Expression, "Customer Address")
	}
	if samePackageDecl.Declaration.Name != "CustomerAddress" {
		t.Fatalf("got declaration %q, want %q", samePackageDecl.Declaration.Name, "CustomerAddress")
	}
	if len(samePackageDecl.Inputs) != 2 {
		t.Fatalf("got %d inputs, want 2", len(samePackageDecl.Inputs))
	}
	if samePackageDecl.Inputs[0].Expr != "Customer" {
		t.Fatalf("got first expr %q, want %q", samePackageDecl.Inputs[0].Expr, "Customer")
	}
	if samePackageDecl.Inputs[0].Type == nil {
		t.Fatal("expected first resolved type")
	}
	if samePackageDecl.Inputs[0].Struct == nil {
		t.Fatal("expected first resolved struct")
	}
	if samePackageDecl.Inputs[1].Expr != "Address" {
		t.Fatalf("got second expr %q, want %q", samePackageDecl.Inputs[1].Expr, "Address")
	}
	if samePackageDecl.Inputs[1].Type == nil {
		t.Fatal("expected second resolved type")
	}
	if samePackageDecl.Inputs[1].Struct == nil {
		t.Fatal("expected second resolved struct")
	}

	generic := resolved[genericIndex]
	if generic.Declaration.Expression != "Customer Envelope[string]" {
		t.Fatalf("got declaration expression %q, want %q", generic.Declaration.Expression, "Customer Envelope[string]")
	}
	if generic.Declaration.Name != "CustomerEnvelope" {
		t.Fatalf("got declaration %q, want %q", generic.Declaration.Name, "CustomerEnvelope")
	}
	if len(generic.Inputs) != 2 {
		t.Fatalf("got %d generic inputs, want 2", len(generic.Inputs))
	}
	if generic.Inputs[0].Expr != "Customer" {
		t.Fatalf("got generic first expr %q, want %q", generic.Inputs[0].Expr, "Customer")
	}
	if generic.Inputs[1].Expr != "Envelope[string]" {
		t.Fatalf("got generic second expr %q, want %q", generic.Inputs[1].Expr, "Envelope[string]")
	}
	if generic.Inputs[1].Type == nil {
		t.Fatal("expected generic resolved type")
	}
	if generic.Inputs[1].Struct == nil {
		t.Fatal("expected generic resolved struct")
	}
	if generic.Inputs[1].Struct.NumFields() != 1 {
		t.Fatalf("got %d generic struct fields, want 1", generic.Inputs[1].Struct.NumFields())
	}
	if generic.Inputs[1].Struct.Field(0).Name() != "Value" {
		t.Fatalf("got generic field %q, want %q", generic.Inputs[1].Struct.Field(0).Name(), "Value")
	}
	if got := types.TypeString(generic.Inputs[1].Struct.Field(0).Type(), nil); got != "string" {
		t.Fatalf("got generic field type %q, want %q", got, "string")
	}

	imported := resolved[importedIndex]
	if imported.Declaration.Expression != "Customer timex.Time" {
		t.Fatalf("got imported declaration expression %q, want %q", imported.Declaration.Expression, "Customer timex.Time")
	}
	if len(imported.Inputs) != 2 {
		t.Fatalf("got %d imported inputs, want 2", len(imported.Inputs))
	}
	if imported.Inputs[1].Expr != "timex.Time" {
		t.Fatalf("got imported expr %q, want %q", imported.Inputs[1].Expr, "timex.Time")
	}
	if imported.Inputs[1].Struct == nil {
		t.Fatal("expected imported resolved struct")
	}

	alias := resolved[aliasIndex]
	if alias.Declaration.Expression != "Customer LocalTime" {
		t.Fatalf("got alias declaration expression %q, want %q", alias.Declaration.Expression, "Customer LocalTime")
	}
	if len(alias.Inputs) != 2 {
		t.Fatalf("got %d alias inputs, want 2", len(alias.Inputs))
	}
	if alias.Inputs[1].Expr != "LocalTime" {
		t.Fatalf("got alias expr %q, want %q", alias.Inputs[1].Expr, "LocalTime")
	}
	if alias.Inputs[1].Struct == nil {
		t.Fatal("expected alias resolved struct")
	}

	typeParam := resolved[typeParamIndex]
	if typeParam.Declaration.Expression != "Customer Envelope[T]" {
		t.Fatalf("got type parameter declaration expression %q, want %q", typeParam.Declaration.Expression, "Customer Envelope[T]")
	}
	if len(typeParam.Inputs) != 2 {
		t.Fatalf("got %d type parameter inputs, want 2", len(typeParam.Inputs))
	}
	if typeParam.Inputs[1].Expr != "Envelope[T]" {
		t.Fatalf("got type parameter expr %q, want %q", typeParam.Inputs[1].Expr, "Envelope[T]")
	}
	if typeParam.Inputs[1].Struct == nil {
		t.Fatal("expected type parameter resolved struct")
	}
	if got := types.TypeString(typeParam.Inputs[1].Struct.Field(0).Type(), nil); got != "T" {
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
