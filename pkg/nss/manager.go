package nss

import (
	"fmt"
	"path/filepath"
	"strings"
)

// databasePattern is the pattern used to find NSS databases with dynamic path configuration.
const databasePattern = "cert*.db"

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

// AppendDefaultNSSDatabasePaths appends the usual NSS database paths to the existing NSS database paths.
func (m *Manager) AppendDefaultNSSDatabasePaths() error {
	theoricPath, err := defaultNSSDatabasePaths()
	if err != nil {
		return err
	}

	dbPaths := m.expandAndFilterPaths(filepath.Glob, theoricPath)

	m.dbPath = append(m.dbPath, dbPaths...)

	return nil
}

// expandAndFilterPaths expands path patterns and filters out directories not containing a database.
func (m *Manager) expandAndFilterPaths(expander func(string) ([]string, error), paths []string) []string {
	var dbPath []string

	for _, path := range paths {
		matches, err := expander(path)
		if err != nil {
			m.Errorf("failed to parse NSS database file pattern (%s):  %v", path, err)

			continue
		}

		// if the path is a pattern, we need to filter the directories not containing a database
		if strings.Contains(path, "*") {
			for _, match := range matches {
				dbFiles, _ := expander(filepath.Join(match, databasePattern))
				if len(dbFiles) > 0 {
					dbPath = append(dbPath, match)
				}
			}
		} else {
			dbPath = append(dbPath, matches...)
		}
	}

	return dbPath
}
