package writer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteFileCreatesGeneratedOutput(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "zz_generated.adtgen.go")

	if err := WriteFile(path, "package sample\n"); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("os.Stat() error = %v", err)
	}
}
