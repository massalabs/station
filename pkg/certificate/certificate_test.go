package certificate

import (
	"testing"
)

func TestReadFile(t *testing.T) {
	t.Run("empty path", func(t *testing.T) {
		_, err := readFile("")
		if err == nil {
			t.Errorf("Expected an error for empty path but got nil")
		}
	})

	t.Run("file not found", func(t *testing.T) {
		_, err := readFile("testdata/non_existent_file.pem")
		if err == nil {
			t.Errorf("Expected an error for non-existent file but got nil")
		}
	})

	t.Run("successful read", func(t *testing.T) {
		_, err := readFile("testdata/cert.pem")
		if err != nil {
			t.Errorf("Expected no error for valid file but got %v", err)
		}
	})
}

func TestLoadCertificate(t *testing.T) {
	t.Run("invalid certificate", func(t *testing.T) {
		// Assuming you have a invalid_cert.pem file
		_, err := LoadCertificate("invalid_cert.pem")
		if err == nil {
			t.Errorf("Expected an error for invalid certificate but got nil")
		}
	})

	t.Run("successful load", func(t *testing.T) {
		// Assuming you have a valid pem file
		_, err := LoadCertificate("testdata/cert.pem")
		if err != nil {
			t.Errorf("Expected no error for valid certificate but got %v", err)
		}
	})
}

func TestLoadPrivateKey(t *testing.T) {
	t.Run("invalid private key", func(t *testing.T) {
		// Assuming you have a invalid_key.pem file
		_, err := LoadPrivateKey("invalid_key.pem")
		if err == nil {
			t.Errorf("Expected an error for invalid private key but got nil")
		}
	})

	t.Run("successful load", func(t *testing.T) {
		// Assuming you have a valid key.pem file
		_, err := LoadPrivateKey("testdata/key.pem")
		if err != nil {
			t.Errorf("Expected no error for valid private key but got %v", err)
		}
	})
}
