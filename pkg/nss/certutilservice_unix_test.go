package nss

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCertUtilService_AddCA tests the AddCA method of CertUtilService.
func TestCertUtilService_AddCA(t *testing.T) {
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		t.Skip("skipping test; CI environment detected")
	}

	dbPath := "testdata"
	certPath := "testdata/cert.pem"
	certName := "testNSS"

	runner, err := NewCertUtilRunner()
	if err != nil {
		t.Fatalf("Failed to create runner: %v", err)
	}

	service, err := NewCertUtilService(runner)
	if err != nil {
		t.Fatalf("Failed to create certutil service: %v", err)
	}

	// Adding the certificate.
	err = service.AddCA(dbPath, certName, certPath)
	assert.Nil(t, err)

	// Verifying the certificate is present.
	err = service.IsKnownCA(dbPath, certName)
	assert.Nil(t, err)

	// Removing the certificate.
	err = service.DeleteCA(dbPath, certName)
	assert.Nil(t, err)

	// Verifying the certificate is no longer present.
	err = service.IsKnownCA(dbPath, certName)
	assert.Error(t, err)
}
