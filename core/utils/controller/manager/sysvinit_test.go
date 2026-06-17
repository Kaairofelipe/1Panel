package manager

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSysvinit_IsActive(t *testing.T) {
	// Create a temporary directory for our mock commands
	tmpDir, err := os.MkdirTemp("", "sysvinit-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a mock 'service' command
	mockService := filepath.Join(tmpDir, "service")
	// The script will use /bin/true or /bin/false via exit status equivalents
	script := "#!/bin/sh\nif [ \"$1\" = \"active_service\" ] && [ \"$2\" = \"status\" ]; then\n\techo 'active' > /dev/null\nelse\n\techo 'inactive' > /dev/null\n\texec /bin/false\nfi\n"

	if err := os.WriteFile(mockService, []byte(script), 0755); err != nil {
		t.Fatalf("failed to write mock service script: %v", err)
	}

	// Prepend the temp directory to PATH so our mock is found first
	oldPath := os.Getenv("PATH")
	if err := os.Setenv("PATH", tmpDir+string(os.PathListSeparator)+oldPath); err != nil {
		t.Fatalf("failed to set PATH: %v", err)
	}
	defer os.Setenv("PATH", oldPath)

	sysvinit := NewSysvinit()

	tests := []struct {
		name        string
		serviceName string
		wantActive  bool
	}{
		{
			name:        "active service",
			serviceName: "active_service",
			wantActive:  true,
		},
		{
			name:        "inactive service",
			serviceName: "inactive_service",
			wantActive:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotActive, err := sysvinit.IsActive(tt.serviceName)
			if err != nil {
				t.Fatalf("IsActive() error = %v, want nil", err)
			}
			if gotActive != tt.wantActive {
				t.Errorf("IsActive() = %v, want %v", gotActive, tt.wantActive)
			}
		})
	}
}
