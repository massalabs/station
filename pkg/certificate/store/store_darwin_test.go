package store

import (
	"os"
	"testing"

	"github.com/massalabs/station/pkg/certificate"
	"github.com/stretchr/testify/require"
)

func TestManualCheck(t *testing.T) {
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		t.Skip("skipping test; CI environment detected")
	}

	if os.Geteuid() != 0 {
		t.Skip("skipping test; not running as root")
	}

	certPath := "../testdata/cert.pem"

	cert, err := certificate.LoadCertificate(certPath)
	if err != nil {
		t.Fatal(err)
	}

	err = Add(cert)
	require.NoError(t, err)
}
