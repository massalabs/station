//go:build darwin

package store

import (
	"bytes"
	"crypto/sha1" //nolint:gosec // SHA-1 is how the keychain identifies certificates, not used for security.
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"strings"

	"github.com/massalabs/station/pkg/logger"
)

const (
	permissionUrwGrOr = 0o644

	systemKeychain  = "/Library/Keychains/System.keychain"
	trustedCertFile = "trusted-cert"
)

// Add trusts the certificate as a root CA for TLS.
// Certificates with the same subject (e.g. a previous, expired CA) are removed first so that
// a renewed CA does not coexist with stale ones.
//
// When running as root (installer), the certificate is added to the system keychain with admin trust settings.
// Otherwise (e.g. CA renewal at runtime), it is added to the user's default (login) keychain with user
// trust settings: changing admin trust settings requires an interactive authorization that elevated
// non-interactive processes cannot get, whereas the user domain triggers the standard macOS password prompt.
func Add(cert *x509.Certificate) error {
	certFile, err := os.CreateTemp("", trustedCertFile)
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}

	defer func() {
		_ = certFile.Close()
		_ = os.Remove(certFile.Name())
	}()

	err = os.WriteFile(certFile.Name(), pem.EncodeToMemory(
		&pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw}), fs.FileMode(permissionUrwGrOr))
	if err != nil {
		return fmt.Errorf("failed to write certificate: %w", err)
	}

	// Explicit ssl and basic policies with trustRoot result.
	// The keychain must always be given: without it, only the trust settings are recorded and the certificate
	// is not imported, so clients receiving only the leaf certificate cannot build the chain to the CA.
	addArgs := []string{"add-trusted-cert", "-r", "trustRoot", "-p", "ssl", "-p", "basic"}
	deleteArgs := []string{}

	var keychain string

	if os.Geteuid() == 0 {
		keychain = systemKeychain
		// Also remove the admin trust settings of the stale certificates.
		deleteArgs = append(deleteArgs, "-t")
		addArgs = append(addArgs, "-d")
	} else {
		keychain, err = defaultKeychain()
		if err != nil {
			return err
		}
	}

	addArgs = append(addArgs, "-k", keychain)

	staleHashes, err := findStaleCertificates(cert, keychain)
	if err != nil {
		// Not blocking: the new certificate can still be trusted.
		logger.Warnf("failed to look for stale certificates in the keychain: %s", err)
	}

	for _, hash := range staleHashes {
		args := append([]string{"delete-certificate", "-Z", hash}, deleteArgs...)

		err = runSecurity(append(args, keychain)...)
		if err != nil {
			logger.Warnf("failed to remove stale certificate %s from the keychain: %s", hash, err)
		}
	}

	err = runSecurity(append(addArgs, certFile.Name())...)
	if err != nil {
		return fmt.Errorf("failed to add the certificate to the keychain: %w", err)
	}

	return nil
}

func Delete(_ *x509.Certificate) error {
	return errors.New("not implemented")
}

// defaultKeychain returns the path of the user's default (login) keychain.
func defaultKeychain() (string, error) {
	out, err := exec.Command("security", "default-keychain").Output()
	if err != nil {
		return "", fmt.Errorf("failed to get the default keychain: %w", err)
	}

	keychain := strings.Trim(strings.TrimSpace(string(out)), `"`)
	if keychain == "" {
		return "", errors.New("no default keychain")
	}

	return keychain, nil
}

// findStaleCertificates returns the SHA-1 hashes of the certificates of the given keychain
// having the same subject as the given certificate but being a different certificate.
func findStaleCertificates(cert *x509.Certificate, keychain string) ([]string, error) {
	//nolint:gosec // arguments are not user controlled.
	out, err := exec.Command(
		"security", "find-certificate", "-a", "-c", cert.Subject.CommonName, "-p", keychain,
	).Output()
	if err != nil {
		// security exits with an error when no certificate matches.
		return nil, nil
	}

	var hashes []string

	for block, rest := pem.Decode(out); block != nil; block, rest = pem.Decode(rest) {
		existing, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			continue
		}

		if !bytes.Equal(existing.RawSubject, cert.RawSubject) || bytes.Equal(existing.Raw, cert.Raw) {
			continue
		}

		sum := sha1.Sum(existing.Raw) //nolint:gosec
		hashes = append(hashes, strings.ToUpper(hex.EncodeToString(sum[:])))
	}

	return hashes, nil
}

// runSecurity runs the security command with the given arguments.
func runSecurity(args ...string) error {
	out, err := exec.Command("security", args...).CombinedOutput() //nolint:gosec // arguments are not user controlled.
	if err != nil {
		return fmt.Errorf("security %s: %s: %w", args[0], strings.TrimSpace(string(out)), err)
	}

	return nil
}
