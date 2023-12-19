package nss

// This file provides a Manager struct that encapsulates operations on all NSS (Network Security Services)
// databases within an operating system. NSS databases typically contain SSL, TLS and other cryptographic certificates.
//
// The Manager struct holds a list of usual OS specific database paths (dbPath), an interface to the CertUtilServicer
// for certutil command execution, and a Logger interface for logging operations.
//
// The Manager struct is primarily used for batch operations on all NSS databases. It uses a CertUtilServicer to execute
// operations on individual databases, and uses the Logger interface to log errors or debug messages.
//
// Future enhancements:
// More functionalities can be added to the Manager struct as needed, based on the required interactions with the NSS
// databases.

import (
	"fmt"
	"path/filepath"
	"strings"
)

// databasePattern is the pattern used to find NSS databases with dynamic path configuration.
const databasePattern = "cert*.db"

// Manager encapsulates operations to execute on all the NSS databases of the operating system.
type Manager struct {
	dbPath   []string
	certutil CertUtilServicer
	Logger
}

// NewManager returns a new Manager instance.
// It will manage the NSS databases corresponding to the given paths.
func NewManager(dbPath []string, certutil CertUtilServicer, logger Logger) *Manager {
	return &Manager{
		dbPath:   dbPath,
		certutil: certutil,
		Logger:   logger,
	}
}

// executeOnPaths executes the given operation on each NSS database path.
func (m *Manager) executeOnPaths(operation func(path string) error) error {
	var chainedErr error

	for _, path := range m.dbPath {
		if err := operation(path); err != nil {
			chainedErr = fmt.Errorf("%w\n%w", chainedErr, err)
		}
	}

	return chainedErr
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
		if err != nil {
			m.Debugf("certificate is unknown by the NSS database (%s): %v", path, err)

			return false
		}

		m.Debugf("certificate is known by the NSS database (%s)", path)
	}

	return true
}

// AppendDefaultNSSDatabasePaths appends the usual NSS database paths to the given NSS database paths.
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
