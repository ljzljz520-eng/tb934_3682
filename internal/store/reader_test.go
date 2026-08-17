package store

import (
	"path/filepath"
	"testing"
)

func TestReaderLimitCanBeReused(t *testing.T) {
	repository, err := NewWithReaderLimit(filepath.Join(t.TempDir(), "readers.db"), 2)
	if err != nil {
		t.Fatal(err)
	}
	defer repository.Close()
	first, err := repository.OpenReader()
	if err != nil {
		t.Fatal(err)
	}
	second, err := repository.OpenReader()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := repository.OpenReader(); err != ErrReaderLimit {
		t.Fatalf("expected reader limit, got %v", err)
	}
	if err := first.Close(); err != nil {
		t.Fatal(err)
	}
	if err := second.Close(); err != nil {
		t.Fatal(err)
	}
	third, err := repository.OpenReader()
	if err != nil {
		t.Fatal(err)
	}
	if err := third.Close(); err != nil {
		t.Fatal(err)
	}
}
