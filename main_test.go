package main

import (
	"path/filepath"
	"testing"
)

func TestOutputPathFromSourceFilename(t *testing.T) {
	filename := filepath.Join("tmp", "example", "generate_types.go")
	got, err := outputPathFromSourceFilename(filename)
	if err != nil {
		t.Fatalf("outputPathFromSourceFilename() error = %v", err)
	}
	want := filepath.Join("tmp", "example", "generate_types_adtgen.go")
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestOutputPathFromSourceFilenameRejectsNonGoFiles(t *testing.T) {
	filename := filepath.Join("tmp", "example", "generate_types.txt")
	_, err := outputPathFromSourceFilename(filename)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
