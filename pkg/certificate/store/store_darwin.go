//go:build darwin

package store

import (
	"bytes"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"io/fs"
	"os"
	"os/exec"

	"github.com/massalabs/station/pkg/runner"
	"howett.net/plist"
)

var _ runner.Runner = &SecurityRunner{}

const (
	permissionUrwGrOr = 0o644
	permissionUrw     = 0o600

	tmpTrustSettingsFile = "trust-settings.plist"
	tmpTrustedCertFile   = "trusted-cert"
)

// SecurityRunner encapsulates security commands.
type SecurityRunner struct {
	runner.CommandRunner
}

// NewSecurityRunner returns a new SecurityRunner.
// It returns an error if the security binary is not found.
func NewSecurityRunner() (*SecurityRunner, error) {
	bin, err := exec.LookPath("security")
	if err != nil {
		return nil, fmt.Errorf("failed to find security binary: %w", err)
	}

	return &SecurityRunner{runner.CommandRunner{BinaryPath: bin}}, nil
}

func Add(cert *x509.Certificate) error {
	security, err := NewSecurityRunner()
	if err != nil {
		return fmt.Errorf("failed to instantiate the certutil runner: %w", err)
	}

	err = addTrustedCert(cert, security)
	if err != nil {
		return fmt.Errorf("failed to add the certificate to the system keychain: %w", err)
	}

	plistRoot, err := exportTrustSettingsContent(security)
	if err != nil {
		return fmt.Errorf("failed to export trust settings: %w", err)
	}

	err = updateTrustSettings(plistRoot, cert)
	if err != nil {
		return fmt.Errorf("failed to update trust settings: %w", err)
	}

	returnValue := importTrustSettings(plistRoot, security)
	if returnValue != nil {
		return fmt.Errorf("failed to re-import trust settings: %w", err)
	}

	return nil
}

func Delete(_ *x509.Certificate) error {
	return fmt.Errorf("not implemented")
}

// addTrustedCert adds the certificate to the system keychain.
func addTrustedCert(cert *x509.Certificate, security *SecurityRunner) error {
	trustedCertFile, err := os.CreateTemp("", tmpTrustedCertFile)
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer trustedCertFile.Close()

	err = os.WriteFile(trustedCertFile.Name(), pem.EncodeToMemory(
		&pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw}), fs.FileMode(permissionUrwGrOr))
	if err != nil {
		return fmt.Errorf("failed to write certificate: %w", err)
	}

	err = security.Run(
		"add-trusted-cert", "-d", "-k", "/Library/Keychains/System.keychain", trustedCertFile.Name(),
	)
	if err != nil {
		return fmt.Errorf("failed to add the certificate to the system keychain: %w", err)
	}

	return nil
}

// exportTrustSettingsContent exports the trust settings to a temporary file.
func exportTrustSettingsContent(security *SecurityRunner) (map[string]interface{}, error) {
	plistFile, err := os.CreateTemp("", tmpTrustSettingsFile)
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary trust-settings file: %w", err)
	}
	defer plistFile.Close()

	err = security.Run("trust-settings-export", "-d", plistFile.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall plistData: %w", err)
	}

	plistData, err := os.ReadFile(plistFile.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall plistData: %w", err)
	}

	var plistRoot map[string]interface{}

	_, err = plist.Unmarshal(plistData, &plistRoot)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall plistData: %w", err)
	}

	return plistRoot, nil
}

// importTrustSettings imports the trust settings from a temporary file.
func importTrustSettings(plistRoot map[string]interface{}, security *SecurityRunner) error {
	updatedPlistData, err := plist.MarshalIndent(plistRoot, plist.XMLFormat, "\t")
	if err != nil {
		return fmt.Errorf("failed to marshal trust settings: %w", err)
	}

	err = os.WriteFile(tmpTrustSettingsFile, updatedPlistData, fs.FileMode(permissionUrw))
	if err != nil {
		return fmt.Errorf("failed to write trust settings: %w", err)
	}

	err = security.Run("trust-settings-import", "-d", tmpTrustSettingsFile)
	if err != nil {
		return fmt.Errorf("failed to re-import settings")
	}
	return nil
}

// updateTrustSettings updates the trust settings with the certificate.
func updateTrustSettings(plistRoot map[string]interface{}, cert *x509.Certificate) error {
	if plistRoot["trustVersion"].(uint64) != 1 {
		return fmt.Errorf("unsupported trust settings version: %d", plistRoot["trustVersion"])
	}

	rootSubjectASN1, err := asn1.Marshal(cert.Subject.ToRDNSequence())
	if err != nil {
		return fmt.Errorf("failed to marshal cert root subject: %w", err)
	}

	trustSettings, err := createCertTrustSettings()
	if err != nil {
		return fmt.Errorf("failed to create trust settings: %w", err)
	}

	trustList := plistRoot["trustList"].(map[string]interface{})

	for key := range trustList {
		entry := trustList[key].(map[string]interface{})

		if _, ok := entry["issuerName"]; !ok {
			continue
		}

		issuerName := entry["issuerName"].([]byte)

		if !bytes.Equal(rootSubjectASN1, issuerName) {
			continue
		}

		entry["trustSettings"] = trustSettings

		break
	}

	return nil
}

// createCertTrustSettings creates the trust settings for the certificate.
func createCertTrustSettings() ([]interface{}, error) {
	var trustSettings []interface{}

	trustSettingsData := []byte(`
<array>
	<dict>
		<key>kSecTrustSettingsPolicy</key>
		<data>
		KoZIhvdjZAED
		</data>
		<key>kSecTrustSettingsPolicyName</key>
		<string>sslServer</string>
		<key>kSecTrustSettingsResult</key>
		<integer>1</integer>
	</dict>
	<dict>
		<key>kSecTrustSettingsPolicy</key>
		<data>
		KoZIhvdjZAEC
		</data>
		<key>kSecTrustSettingsPolicyName</key>
		<string>basicX509</string>
		<key>kSecTrustSettingsResult</key>
		<integer>1</integer>
	</dict>
</array>
`)

	_, err := plist.Unmarshal(trustSettingsData, &trustSettings)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal trust settings: %w", err)
	}

	return trustSettings, nil
}
