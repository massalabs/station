package store

import (
	"os"
	"testing"

	"github.com/massalabs/station/pkg/certificate"
	"github.com/stretchr/testify/assert"
)

func TestManualCheck(t *testing.T) {
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		t.Skip("skipping test; CI environment detected")
	}
	certPath := "../testdata/cert.pem"

	cert, err := certificate.LoadCertificate(certPath)
	if err != nil {
		t.Fatal(err)
	}

	err = Add(cert)
	assert.NoError(t, err)
}
