package manager

import (
	"errors"
	"testing"

	"github.com/1Panel-dev/1Panel/core/utils/ssh"
)

func TestSystemd_IsEnable(t *testing.T) {
	tests := []struct {
		name        string
		serviceName string
		mockRun     func(client *ssh.SSHClient, name string, args ...string) (string, error)
		want        bool
		wantErr     bool
	}{
		{
			name:        "enabled service",
			serviceName: "nginx",
			mockRun: func(client *ssh.SSHClient, name string, args ...string) (string, error) {
				return "enabled\n", nil
			},
			want:    true,
			wantErr: false,
		},
		{
			name:        "disabled service",
			serviceName: "nginx",
			mockRun: func(client *ssh.SSHClient, name string, args ...string) (string, error) {
				return "disabled\n", errors.New("exit status 1")
			},
			want:    false,
			wantErr: false,
		},
		{
			name:        "sshd alias",
			serviceName: "sshd",
			mockRun: func(client *ssh.SSHClient, name string, args ...string) (string, error) {
				if len(args) > 0 && args[len(args)-1] == "sshd" {
					return "alias\n", errors.New("exit status 1")
				}
				if len(args) > 0 && args[len(args)-1] == "ssh" {
					return "enabled\n", nil
				}
				return "", errors.New("unknown")
			},
			want:    true,
			wantErr: false,
		},
		{
			name:        "snap fallback enabled",
			serviceName: "docker",
			mockRun: func(client *ssh.SSHClient, name string, args ...string) (string, error) {
				if name == "systemctl" {
					return "unknown\n", errors.New("exit status 1")
				}
				if name == "snap" && len(args) > 0 && args[0] == "services" {
					return "Service  Startup  Current  Notes\ndocker   enabled  active   -", nil
				}
				return "", errors.New("unknown")
			},
			want:    true,
			wantErr: false,
		},
		{
			name:        "error case",
			serviceName: "nginx",
			mockRun: func(client *ssh.SSHClient, name string, args ...string) (string, error) {
				if name == "systemctl" {
					return "error", errors.New("systemctl error")
				}
				if name == "snap" {
					return "error", errors.New("snap error")
				}
				return "", errors.New("unknown")
			},
			want:    false,
			wantErr: true,
		},
	}

	originalRun := run
	defer func() { run = originalRun }()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			run = tt.mockRun
			s := NewSystemd()
			got, err := s.IsEnable(tt.serviceName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Systemd.IsEnable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Systemd.IsEnable() = %v, want %v", got, tt.want)
			}
		})
	}
}
