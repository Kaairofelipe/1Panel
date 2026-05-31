package files

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileInfoSearch(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "searchtest")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	filesToCreate := []string{
		"test_file1.txt",
		"test_file2.log",
		"other_file.txt",
		"nested/dir/test_file3.txt",
	}

	for _, f := range filesToCreate {
		fullPath := filepath.Join(tempDir, f)
		err := os.MkdirAll(filepath.Dir(fullPath), 0755)
		if err != nil {
			t.Fatalf("failed to create dir: %v", err)
		}
		err = os.WriteFile(fullPath, []byte("test"), 0644)
		if err != nil {
			t.Fatalf("failed to create file: %v", err)
		}
	}

	fi := &FileInfo{
		Path: tempDir,
	}

	results, total, err := fi.search("test_file", 10)
	if err != nil {
		t.Fatalf("search failed: %v", err)
	}
	if total != 3 {
		t.Errorf("expected 3 total files, got %d", total)
	}
	if len(results) != 3 {
		t.Errorf("expected 3 results, got %d", len(results))
	}

	results, total, err = fi.search("test_file", 2)
	if err != nil {
		t.Fatalf("search failed: %v", err)
	}
	if total != 3 {
		t.Errorf("expected 3 total files, got %d", total)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
}
