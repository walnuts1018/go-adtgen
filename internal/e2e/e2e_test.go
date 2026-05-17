package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestGoGenerateProducesExpectedOutput(t *testing.T) {
	dir := filepath.Join("..", "testdata", "fixtures", "e2e", "base")
	outputPath := filepath.Join(dir, "zz_generated.product_types.go")
	_ = os.Remove(outputPath)

	cmd := exec.Command("go", "generate", ".")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go generate failed: %v\n%s", err, out)
	}

	got, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("os.ReadFile(got) error = %v", err)
	}
	want, err := os.ReadFile(filepath.Join(dir, "expected.txt"))
	if err != nil {
		t.Fatalf("os.ReadFile(want) error = %v", err)
	}
	if string(got) != string(want) {
		t.Fatalf("generated output mismatch\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

func TestGeneratedFixtureBuilds(t *testing.T) {
	cmd := exec.Command("go", "test", "../testdata/fixtures/e2e/base")
	cmd.Dir = "."
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("fixture build failed: %v\n%s", err, out)
	}
}
