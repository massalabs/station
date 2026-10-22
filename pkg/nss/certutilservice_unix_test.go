package nss

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestCertUtilService_AddCA tests the AddCA method of CertUtilService.
func TestCertUtilService_AddCA(t *testing.T) {
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		t.Skip("skipping test; CI environment detected")
	}

	// The NSS database and the certificate are generated for each run,
	// so the test neither expires nor modifies tracked files.
	dbPath := t.TempDir()
	certPath := writeSelfSignedCA(t, t.TempDir())
	certName := "testNSS"

	runner, err := NewCertUtilRunner()
	if err != nil {
		t.Fatalf("Failed to create runner: %v", err)
	}

	// Creating an empty NSS database.
	err = runner.Run("-N", "--empty-password", "-d", dbPath)
	require.NoError(t, err)

	service, err := NewCertUtilService(runner)
	if err != nil {
		t.Fatalf("Failed to create certutil service: %v", err)
	}

	// Adding the certificate.
	err = service.AddCA(dbPath, certName, certPath)
	require.NoError(t, err)

	// Verifying the certificate is present.
	err = service.IsKnownCA(dbPath, certName)
	require.NoError(t, err)

	// Removing the certificate.
	err = service.DeleteCA(dbPath, certName)
	require.NoError(t, err)

	// Verifying the certificate is no longer present.
	err = service.IsKnownCA(dbPath, certName)
	require.Error(t, err)
}

// writeSelfSignedCA writes a short-lived self-signed CA certificate in PEM format to dir and returns its path.
func writeSelfSignedCA(t *testing.T, dir string) string {
	t.Helper()

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	template := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{Organization: []string{"Massa Labs"}, CommonName: "testing nss"},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(24 * time.Hour),
		IsCA:                  true,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
	}

	der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	require.NoError(t, err)

	certPath := filepath.Join(dir, "cert.pem")
	err = os.WriteFile(certPath, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}), 0o600)
	require.NoError(t, err)

	return certPath
}
