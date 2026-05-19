package e2e

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestGoGenerateProducesExpectedOutput(t *testing.T) {
	dir := filepath.Join("..", "testdata", "fixtures", "e2e", "base")
	outputPath := filepath.Join(dir, "generate_types_adtgen.go")
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

func TestGoGenerateBuildsAndExercisesSumFixture(t *testing.T) {
	dir := filepath.Join("..", "testdata", "fixtures", "e2e", "sum")
	outputPath := filepath.Join(dir, "generate_types_adtgen.go")
	_ = os.Remove(outputPath)

	cmd := exec.Command("go", "generate", ".")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go generate failed: %v\n%s", err, out)
	}

	cmd = exec.Command("go", "test", ".")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("sum fixture test failed: %v\n%s", err, out)
	}
}

func TestGoGenerateSeparatesOutputsPerSourceFile(t *testing.T) {
	dir := filepath.Join("..", "testdata", "fixtures", "e2e", "multi")
	alphaOutput := filepath.Join(dir, "generate_alpha_adtgen.go")
	betaOutput := filepath.Join(dir, "generate_beta_adtgen.go")
	_ = os.Remove(alphaOutput)
	_ = os.Remove(betaOutput)

	cmd := exec.Command("go", "generate", ".")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go generate failed: %v\n%s", err, out)
	}

	alpha, err := os.ReadFile(alphaOutput)
	if err != nil {
		t.Fatalf("os.ReadFile(alpha) error = %v", err)
	}
	beta, err := os.ReadFile(betaOutput)
	if err != nil {
		t.Fatalf("os.ReadFile(beta) error = %v", err)
	}

	if !bytes.Contains(alpha, []byte("type Alpha interface")) {
		t.Fatalf("alpha output missing Alpha interface:\n%s", alpha)
	}
	if !bytes.Contains(alpha, []byte("func MatchAlpha")) {
		t.Fatalf("alpha output missing MatchAlpha:\n%s", alpha)
	}
	if !bytes.Contains(beta, []byte("type Beta struct")) {
		t.Fatalf("beta output missing Beta struct:\n%s", beta)
	}
	if !bytes.Contains(beta, []byte("func NewBeta")) {
		t.Fatalf("beta output missing NewBeta:\n%s", beta)
	}
	if bytes.Contains(alpha, []byte("type Beta")) {
		t.Fatalf("alpha output unexpectedly contains beta declaration:\n%s", alpha)
	}
	if bytes.Contains(beta, []byte("type Alpha")) {
		t.Fatalf("beta output unexpectedly contains alpha declaration:\n%s", beta)
	}

	cmd = exec.Command("go", "test", ".")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("multi fixture test failed: %v\n%s", err, out)
	}
}
