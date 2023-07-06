package nss

import (
	"fmt"
)

// manager encapsulates operations to execute on all the NSS databases of the operating system.
type Manager struct {
	dbPath   []string
	certutil CertUtilServicer
	Logger
}

// NewManager returns a new Manager instance.
func NewManager(dbPath []string, certutil CertUtilServicer, logger Logger) *Manager {
	return &Manager{
		dbPath:   dbPath,
		certutil: certutil,
		Logger:   logger,
	}
}

// executeOnPaths executes the given operation on each NSS database path.
func (m *Manager) executeOnPaths(operation func(path string) error) error {
	for _, path := range m.dbPath {
		if err := operation(path); err != nil {
			return err
		}
	}

	return nil
}

// AddCA adds the certificate to the NSS databases.
func (m *Manager) AddCA(certName string, certPath string) error {
	return m.executeOnPaths(func(path string) error {
		err := m.certutil.AddCA(path, certName, certPath)
		m.Debugf("adding the certificate to the NSS database (%s): %s", path, err)

		if err != nil {
			return fmt.Errorf("adding the certificate to the NSS database (%s): %w", path, err)
		}

		return nil
	})
}

// DeleteCA deletes the certificate from the NSS databases.
func (m *Manager) DeleteCA(certName string) error {
	return m.executeOnPaths(func(path string) error {
		err := m.certutil.DeleteCA(path, certName)
		m.Debugf("deleting the certificate from the NSS database (%s): %s", path, err)

		if err != nil {
			return fmt.Errorf("deleting the certificate from the NSS database (%s): %w", path, err)
		}

		return nil
	})
}

// HasCA checks if a certificate is known by the NSS databases.
func (m *Manager) HasCA(certName string) bool {
	for _, path := range m.dbPath {
		err := m.certutil.IsKnownCA(path, certName)
		m.Debugf("checking if the certificate is known by the NSS database (%s): %s", path, err)

		if err != nil {
			return false
		}
	}

	return true
}
