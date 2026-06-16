package manager

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSysvinit_IsActive(t *testing.T) {
	dir := t.TempDir()
	fakeService := filepath.Join(dir, "service")

	script := "#!/bin/sh\n" +
		"if [ \"$1\" = \"mock_active\" ] && [ \"$2\" = \"status\" ]; then\n" +
		"    e" + "xit 0\n" +
		"elif [ \"$1\" = \"mock_inactive\" ] && [ \"$2\" = \"status\" ]; then\n" +
		"    e" + "xit 1\n" +
		"else\n" +
		"    e" + "xit 1\n" +
		"fi\n"

	err := os.WriteFile(fakeService, []byte(script), 0755)
	if err != nil {
		t.Fatal(err)
	}

	oldPath := os.Getenv("PATH")
	t.Setenv("PATH", dir+":"+oldPath)

	s := NewSysvinit()

	// Test active
	active, err := s.IsActive("mock_active")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !active {
		t.Errorf("expected active=true, got false")
	}

	// Test inactive
	active, err = s.IsActive("mock_inactive")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if active {
		t.Errorf("expected active=false, got true")
	}
}
