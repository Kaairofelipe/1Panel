package controller

import (
	"fmt"
	"testing"
)

type MockController struct {
	name   string
	exists map[string]bool
}

func (m *MockController) Name() string {
	return m.name
}

func (m *MockController) IsActive(serviceName string) (bool, error) {
	return false, nil
}

func (m *MockController) IsEnable(serviceName string) (bool, error) {
	return false, nil
}

func (m *MockController) IsExist(serviceName string) (bool, error) {
	exists, ok := m.exists[serviceName]
	return ok && exists, nil
}

func (m *MockController) Status(serviceName string) (string, error) {
	return "", nil
}

func (m *MockController) Operate(operate, serviceName string) error {
	return nil
}

func (m *MockController) Reload() error {
	return nil
}

func TestLoadServiceName(t *testing.T) {
	// Backup and restore newController
	originalNewController := newController
	defer func() { newController = originalNewController }()

	tests := []struct {
		name          string
		mockCtrl      *MockController
		keyword       string
		expected      string
		expectedError bool
	}{
		{
			name: "Systemd - exact match exists",
			mockCtrl: &MockController{
				name:   "systemd",
				exists: map[string]bool{"myservice.service": true},
			},
			keyword:       "myservice",
			expected:      "myservice.service",
			expectedError: false,
		},
		{
			name: "Systemd - predefined match clam",
			mockCtrl: &MockController{
				name:   "systemd",
				exists: map[string]bool{"clamd@scan.service": true},
			},
			keyword:       "clam",
			expected:      "clamd@scan.service",
			expectedError: false,
		},
		{
			name: "Non-systemd - exact match",
			mockCtrl: &MockController{
				name:   "openrc",
				exists: map[string]bool{"myservice": true},
			},
			keyword:       "myservice.service",
			expected:      "myservice",
			expectedError: false,
		},
		{
			name: "No match found",
			mockCtrl: &MockController{
				name:   "systemd",
				exists: map[string]bool{},
			},
			keyword:       "unknown",
			expected:      "",
			expectedError: true,
		},
		{
			name: "Systemd - predefined match fail2ban",
			mockCtrl: &MockController{
				name:   "systemd",
				exists: map[string]bool{"fail2ban": true},
			},
			keyword:       "fail2ban",
			expected:      "fail2ban",
			expectedError: false,
		},
		{
			name: "Systemd - predefined ssh",
			mockCtrl: &MockController{
				name:   "systemd",
				exists: map[string]bool{"sshd": true},
			},
			keyword:       "ssh",
			expected:      "sshd",
			expectedError: false,
		},
		{
			name: "Systemd - handle .service.socket",
			mockCtrl: &MockController{
				name:   "systemd",
				exists: map[string]bool{"myservice.socket": true},
			},
			keyword:       "myservice.service.socket",
			expected:      "myservice.socket",
			expectedError: false,
		},
		{
			name: "NewController returns error",
			mockCtrl: nil,
			keyword: "test",
			expected: "",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newController = func() (Controller, error) {
				if tt.mockCtrl == nil {
					return nil, fmt.Errorf("mock error")
				}
				return tt.mockCtrl, nil
			}

			result, err := LoadServiceName(tt.keyword)

			if tt.expectedError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestLoadProcessedName(t *testing.T) {
	tests := []struct {
		mgr      string
		keyword  string
		expected string
	}{
		{"systemd", "test.service.socket", "test.socket"},
		{"systemd", "test", "test.service"},
		{"systemd", "test.service", "test.service"},
		{"systemd", "test.socket", "test.socket"},
		{"openrc", "test.service.socket", "test.socket"},
		{"openrc", "test.service", "test"},
		{"sysvinit", "test.service", "test"},
		{"openrc", "test", "test"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s-%s", tt.mgr, tt.keyword), func(t *testing.T) {
			result := loadProcessedName(tt.mgr, tt.keyword)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}
