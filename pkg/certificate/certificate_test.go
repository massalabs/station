package certificate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadFile(t *testing.T) {
	tests := []struct {
		name     string
		filepath string
		wantErr  bool
	}{
		{"empty path", "", true},
		{"file not found", "testdata/non_existent_file.pem", true},
		{"successful read", "testdata/cert.pem", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := readFile(tt.filepath)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestLoadCertificate(t *testing.T) {
	tests := []struct {
		name     string
		filepath string
		wantErr  bool
	}{
		{"file not found", "testdata/non_existent_file.pem", true},
		{"invalid PEM type", "testdata/key.pem", true},
		{"invalid certificate", "testdata/invalid_cert.pem", true},
		{"successful load", "testdata/cert.pem", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := LoadCertificate(tt.filepath)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestLoadPrivateKey(t *testing.T) {
	tests := []struct {
		name     string
		filepath string
		wantErr  bool
	}{
		{"file not found", "testdata/non_existent_file.pem", true},
		{"invalid PEM type", "testdata/cert.pem", true},
		{"invalid private key", "testdata/invalid_key.pem", true},
		{"successful load", "testdata/key.pem", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := LoadPrivateKey(tt.filepath)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
