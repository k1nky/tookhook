package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		osargs  []string
		env     map[string]string
		want    Config
		wantErr bool
	}{
		{
			name:   "Default",
			osargs: []string{"tookhook"},
			env:    map[string]string{},
			want: Config{
				Listen:      "localhost:8080",
				DarabaseURI: "hooks.yml",
				LogLevel:    "info",
				Plugins:     "",
				QueueURI:    "127.0.0.1:6379",
			},
			wantErr: false,
		},
		{
			name:   "Only arguments",
			osargs: []string{"tookhook", "-s", "localhost:8000", "-d", "my_hooks.yml", "-l", "debug"},
			env:    map[string]string{},
			want: Config{
				Listen:      "localhost:8000",
				DarabaseURI: "my_hooks.yml",
				LogLevel:    "debug",
				Plugins:     "",
				QueueURI:    "127.0.0.1:6379",
			},
			wantErr: false,
		},
		{
			name:   "Only environment",
			osargs: []string{"tookhook"},
			env: map[string]string{
				"TOOKHOOK_LISTEN":       "localhost:8000",
				"TOOKHOOK_DATABASE_URI": "hooks.yml",
				"TOOKHOOK_LOG_LEVEL":    "debug",
			},
			want: Config{
				Listen:      "localhost:8000",
				DarabaseURI: "hooks.yml",
				LogLevel:    "debug",
				QueueURI:    "127.0.0.1:6379",
			},
			wantErr: false,
		},
		{
			name:   "Environment overrides arguments",
			osargs: []string{"tookhhok", "-s", "localhost:8000", "-d", "hooks.yml", "-l", "debug"},
			env: map[string]string{
				"TOOKHOOK_LISTEN":       "localhost:8080",
				"TOOKHOOK_DATABASE_URI": "my_hooks.yml",
				"TOOKHOOK_LOG_LEVEL":    "warn",
			},
			want: Config{
				Listen:      "localhost:8080",
				DarabaseURI: "my_hooks.yml",
				QueueURI:    "127.0.0.1:6379",
				LogLevel:    "warn",
			},
			wantErr: false,
		},
		{
			name:   "Environment and arguments",
			osargs: []string{"tookhhok", "-s", "localhost:8000"},
			env: map[string]string{
				"TOOKHOOK_LISTEN":       "localhost:8080",
				"TOOKHOOK_DATABASE_URI": "hooks.yml",
			},
			want: Config{
				Listen:      "localhost:8080",
				DarabaseURI: "hooks.yml",
				QueueURI:    "127.0.0.1:6379",
				LogLevel:    "info",
			},
			wantErr: false,
		},
		{
			name:    "With invalid argument",
			osargs:  []string{"tookhhok", "-t"},
			env:     map[string]string{},
			want:    Config{},
			wantErr: true,
		},
		{
			name:    "With invalid argument value",
			osargs:  []string{"tookhhok", "-s", "127.0.0.1/8000"},
			env:     map[string]string{},
			want:    Config{},
			wantErr: true,
		},
		{
			name:    "With invalid evironment variable value",
			osargs:  []string{"tookhhok"},
			env:     map[string]string{"TOOKHOOK_LISTEN": "127.0.0.1/8000"},
			want:    Config{},
			wantErr: true,
		},
		{
			name:    "With invalid evironment variable and argument value",
			osargs:  []string{"tookhook", "-s", "127.0.0.2/8000"},
			env:     map[string]string{"TOOKHOOK_LISTEN": "127.0.0.1/8000"},
			want:    Config{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = tt.osargs
			for k, v := range tt.env {
				t.Setenv(k, v)
			}

			c := Config{}
			if err := Parse(&c); err != nil {
				if (err != nil) != tt.wantErr {
					t.Errorf("parseFlags() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			assert.Equal(t, tt.want, c)
		})
	}
}
