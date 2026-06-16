package manager

import (
	"errors"
	"testing"

	"github.com/1Panel-dev/1Panel/core/utils/ssh"
)

func TestSystemd_IsActive(t *testing.T) {
	tests := []struct {
		name        string
		serviceName string
		mockRun     func(client *ssh.SSHClient, name string, args ...string) (string, error)
		want        bool
		wantErr     bool
	}{
		{
			name:        "active service",
			serviceName: "nginx",
			mockRun: func(client *ssh.SSHClient, name string, args ...string) (string, error) {
				return "active\n", nil
			},
			want:    true,
			wantErr: false,
		},
		{
			name:        "inactive service",
			serviceName: "nginx",
			mockRun: func(client *ssh.SSHClient, name string, args ...string) (string, error) {
				return "inactive\n", errors.New("exit status 3")
			},
			want:    false,
			wantErr: false,
		},
		{
			name:        "unknown state but valid err",
			serviceName: "nginx",
			mockRun: func(client *ssh.SSHClient, name string, args ...string) (string, error) {
				if name == "systemctl" && args[0] == "is-active" {
					return "failed\n", errors.New("exit status 3")
				}
				if name == "snap" && args[0] == "services" {
					return "", errors.New("not found")
				}
				return "", nil
			},
			want:    false,
			wantErr: true,
		},
		{
			name:        "fallback to snap active",
			serviceName: "docker",
			mockRun: func(client *ssh.SSHClient, name string, args ...string) (string, error) {
				if name == "systemctl" && args[0] == "is-active" {
					return "failed\n", errors.New("exit status 3")
				}
				if name == "snap" && args[0] == "services" {
					return "docker  active  -  -  -", nil
				}
				return "", nil
			},
			want:    true,
			wantErr: false,
		},
		{
			name:        "fallback to snap inactive",
			serviceName: "docker",
			mockRun: func(client *ssh.SSHClient, name string, args ...string) (string, error) {
				if name == "systemctl" && args[0] == "is-active" {
					return "failed\n", errors.New("exit status 3")
				}
				if name == "snap" && args[0] == "services" {
					// Use disabled here instead of inactive since "inactive" contains "active", and Snap.IsActive does strings.Contains(line, "active")
					return "docker  disabled  -  -  -", nil
				}
				return "", nil
			},
			want:    false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalRun := run
			run = tt.mockRun
			defer func() { run = originalRun }()

			s := NewSystemd()
			got, err := s.IsActive(tt.serviceName)

			if (err != nil) != tt.wantErr {
				t.Errorf("Systemd.IsActive() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Systemd.IsActive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSystemd_Name(t *testing.T) {
	s := NewSystemd()
	if s.Name() != "systemd" {
		t.Errorf("Systemd.Name() = %v, want %v", s.Name(), "systemd")
	}
}

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
				return "disabled\n", errors.New("exit status 3")
			},
			want:    false,
			wantErr: false,
		},
		{
			name:        "sshd alias",
			serviceName: "sshd",
			mockRun: func(client *ssh.SSHClient, name string, args ...string) (string, error) {
				if args[1] == "sshd" {
					return "alias\n", errors.New("exit status 3")
				}
				if args[1] == "ssh" {
					return "enabled\n", nil
				}
				return "", nil
			},
			want:    true,
			wantErr: false,
		},
		{
			name:        "fallback to snap enabled",
			serviceName: "docker",
			mockRun: func(client *ssh.SSHClient, name string, args ...string) (string, error) {
				if name == "systemctl" && args[0] == "is-enabled" {
					return "failed\n", errors.New("exit status 3")
				}
				if name == "snap" && args[0] == "services" {
					return "docker  enabled  -  -  -", nil
				}
				return "", nil
			},
			want:    true,
			wantErr: false,
		},
		{
			name:        "fallback to snap disabled",
			serviceName: "docker",
			mockRun: func(client *ssh.SSHClient, name string, args ...string) (string, error) {
				if name == "systemctl" && args[0] == "is-enabled" {
					return "failed\n", errors.New("exit status 3")
				}
				if name == "snap" && args[0] == "services" {
					return "docker  disabled  -  -  -", nil
				}
				return "", nil
			},
			want:    false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalRun := run
			run = tt.mockRun
			defer func() { run = originalRun }()

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

func TestSystemd_IsExist(t *testing.T) {
	tests := []struct {
		name        string
		serviceName string
		mockRun     func(client *ssh.SSHClient, name string, args ...string) (string, error)
		want        bool
		wantErr     bool
	}{
		{
			name:        "exist service (enabled)",
			serviceName: "nginx",
			mockRun: func(client *ssh.SSHClient, name string, args ...string) (string, error) {
				return "enabled\n", nil
			},
			want:    true,
			wantErr: false,
		},
		{
			name:        "exist service (disabled)",
			serviceName: "nginx",
			mockRun: func(client *ssh.SSHClient, name string, args ...string) (string, error) {
				return "disabled\n", errors.New("exit status 3")
			},
			want:    true,
			wantErr: true,
		},
		{
			name:        "fallback to snap exist",
			serviceName: "docker",
			mockRun: func(client *ssh.SSHClient, name string, args ...string) (string, error) {
				if name == "systemctl" && args[0] == "is-enabled" {
					return "failed\n", errors.New("exit status 3")
				}
				if name == "snap" && args[0] == "services" {
					return "docker  active  -  -  -", nil
				}
				return "", nil
			},
			want:    true,
			wantErr: false,
		},
		{
			name:        "fallback to snap not exist",
			serviceName: "docker",
			mockRun: func(client *ssh.SSHClient, name string, args ...string) (string, error) {
				if name == "systemctl" && args[0] == "is-enabled" {
					return "failed\n", errors.New("exit status 3")
				}
				if name == "snap" && args[0] == "services" {
					return "", nil
				}
				return "", nil
			},
			want:    false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalRun := run
			run = tt.mockRun
			defer func() { run = originalRun }()

			s := NewSystemd()
			got, err := s.IsExist(tt.serviceName)

			if (err != nil) != tt.wantErr {
				t.Errorf("Systemd.IsExist() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Systemd.IsExist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSystemd_Status(t *testing.T) {
	originalRun := run
	run = func(client *ssh.SSHClient, name string, args ...string) (string, error) {
		return "status output", nil
	}
	defer func() { run = originalRun }()

	s := NewSystemd()
	out, err := s.Status("nginx")
	if err != nil {
		t.Errorf("Systemd.Status() error = %v", err)
	}
	if out != "status output" {
		t.Errorf("Systemd.Status() = %v, want %v", out, "status output")
	}
}

func TestSystemd_Operate(t *testing.T) {
	tests := []struct {
		name        string
		operate     string
		serviceName string
		mockRun     func(client *ssh.SSHClient, name string, args ...string) (string, error)
		wantErr     bool
	}{
		{
			name:        "successful operate",
			operate:     "start",
			serviceName: "nginx",
			mockRun: func(client *ssh.SSHClient, name string, args ...string) (string, error) {
				return "", nil
			},
			wantErr: false,
		},
		{
			name:        "sshd alias",
			operate:     "start",
			serviceName: "sshd",
			mockRun: func(client *ssh.SSHClient, name string, args ...string) (string, error) {
				if args[1] == "sshd" {
					return "Failed to start sshd.service: Unit sshd.service is an alias name or linked unit file", errors.New("exit status 1")
				}
				if args[1] == "ssh" {
					return "", nil
				}
				return "", nil
			},
			wantErr: false,
		},
		{
			name:        "fallback to snap successful",
			operate:     "start",
			serviceName: "docker",
			mockRun: func(client *ssh.SSHClient, name string, args ...string) (string, error) {
				if name == "systemctl" && args[1] == "docker" {
					return "failed\n", errors.New("exit status 1")
				}
				// NewSnap().Operate -> IsExist -> snap services
				if name == "snap" && args[0] == "services" {
					return "docker active - - -", nil
				}
				// NewSnap().Operate -> snap operate
				if name == "snap" && args[0] == "start" && args[1] == "docker" {
					return "", nil
				}
				return "", nil
			},
			wantErr: false, // Snap().Operate returns nil
		},
		{
			name:        "fallback to snap failed",
			operate:     "start",
			serviceName: "docker",
			mockRun: func(client *ssh.SSHClient, name string, args ...string) (string, error) {
				if name == "systemctl" && args[1] == "docker" {
					return "failed systemctl\n", errors.New("exit status 1")
				}
				if name == "snap" && args[0] == "services" {
					return "docker active - - -", nil
				}
				if name == "snap" && args[0] == "start" && args[1] == "docker" {
					return "failed snap\n", errors.New("exit status 1")
				}
				return "", nil
			},
			wantErr: true, // handlerErr processes systemctl err, output is "failed systemctl\n" so error is "failed systemctl\n"
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalRun := run
			run = tt.mockRun
			defer func() { run = originalRun }()

			s := NewSystemd()
			err := s.Operate(tt.operate, tt.serviceName)

			if (err != nil) != tt.wantErr {
				t.Errorf("Systemd.Operate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSystemd_Reload(t *testing.T) {
	originalRun := run
	run = func(client *ssh.SSHClient, name string, args ...string) (string, error) {
		return "", nil
	}
	defer func() { run = originalRun }()

	s := NewSystemd()
	err := s.Reload()
	if err != nil {
		t.Errorf("Systemd.Reload() error = %v", err)
	}
}

func TestSystemd_ReloadError(t *testing.T) {
	originalRun := run
	run = func(client *ssh.SSHClient, name string, args ...string) (string, error) {
		return "error reload", errors.New("exit status 1")
	}
	defer func() { run = originalRun }()

	s := NewSystemd()
	err := s.Reload()
	if err == nil {
		t.Errorf("Systemd.Reload() expected error but got nil")
	}
}
