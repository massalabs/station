package logger

import (
	"testing"
)

func TestLogLevel(t *testing.T) {
	tests := []struct {
		name    string
		env     string
		want    string
		wantErr bool
	}{
		{"Default", "", "INFO", false},
		{"Debug", "DEBUG", "DEBUG", false},
		{"Invalid", "INVALID", "", true},
	}

	//nolint:varnamelen
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("LOG_LEVEL", tt.env)

			got, err := LogLevel()
			if (err != nil) != tt.wantErr {
				t.Errorf("LogLevel() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if got != tt.want {
				t.Errorf("LogLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}
