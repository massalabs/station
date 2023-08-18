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

	"github.com/massalabs/station/pkg/logger"
	"github.com/massalabs/station/pkg/runner"
	"howett.net/plist"
)

var _ runner.Runner = &SecurityRunner{}

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

	tempFile, err := os.CreateTemp("", "cert")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer tempFile.Close()

	permissionUrwGrOr := 0o644

	err = os.WriteFile(tempFile.Name(), pem.EncodeToMemory(
		&pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw}), fs.FileMode(permissionUrwGrOr))
	if err != nil {
		return fmt.Errorf("failed to write certificate: %w", err)
	}

	err = security.Run(
		"add-trusted-cert", "-d", "-k", "/Library/Keychains/System.keychain", tempFile.Name(),
	)
	if err != nil {
		return fmt.Errorf("failed to add the certificate to the system keychain: %w", err)
	}

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

	_, err = plist.Unmarshal(trustSettingsData, &trustSettings)
	if err != nil {
		return fmt.Errorf("failed to unmarshal trust settings: %w", err)
	}

	plistFile, err := os.CreateTemp("", "trust-settings")
	if err != nil {
		return fmt.Errorf("failed to create temporary trust-settings file: %w", err)
	}
	defer tempFile.Close()

	err = security.Run("trust-settings-export", "-d", plistFile.Name())
	if err != nil {
		return fmt.Errorf("failed to export trust settings: %w", err)
	}

	plistData, err := os.ReadFile(plistFile.Name())
	if err != nil {
		return fmt.Errorf("failed to read trust settings: %w", err)
	}

	var plistRoot map[string]interface{}

	_, err = plist.Unmarshal(plistData, &plistRoot)
	if err != nil {
		return fmt.Errorf("failed to unmarshall plistData: %w", err)
	}

	rootSubjectASN1, _ := asn1.Marshal(cert.Subject.ToRDNSequence())

	if plistRoot["trustVersion"].(uint64) != 1 {
		return fmt.Errorf("unsupported trust settings version:", plistRoot["trustVersion"])
	}

	trustList := plistRoot["trustList"].(map[string]interface{})

	for key := range trustList {
		logger.Debug(key)
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

	plistData, err = plist.MarshalIndent(plistRoot, plist.XMLFormat, "\t")
	if err != nil {
		return fmt.Errorf("failed to marshal trust settings: %w", err)
	}

	permissionUrw := 0o600

	err = os.WriteFile(plistFile.Name(), plistData, fs.FileMode(permissionUrw))
	if err != nil {
		return fmt.Errorf("failed to write trust settings: %w", err)
	}

	err = security.Run("trust-settings-import", "-d", plistFile.Name())
	if err != nil {
		return fmt.Errorf("failed to re-import settings")
	}

	return nil
}

func Delete(_ *x509.Certificate) error {
	return fmt.Errorf("not implemented")
}
