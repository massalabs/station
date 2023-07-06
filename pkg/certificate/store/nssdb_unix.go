//go:build unix

package store

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/massalabs/station/pkg/logger"
)

const databasePattern = "cert*.db"

// NSSDatabases is a struct that represents the known NSS databases on the system.
type NSSDatabases struct {
	Paths []string
}

// NewNssDatabases returns a new NSSDatabases instance.
func NewNssDatabases() (*NSSDatabases, error) {
	theoricPath, err := NSSDBPaths()
	if err != nil {
		return nil, err
	}

	return &NSSDatabases{
		Paths: filterExistingPath(theoricPath),
	}, nil
}

// executeOnPaths executes the given operation on each NSS database path.
func (n *NSSDatabases) executeOnPaths(operation func(path string) error) error {
	for _, path := range n.Paths {
		if err := operation(path); err != nil {
			return err
		}
	}

	return nil
}

// Add adds the certificate to the NSS databases.
func (n *NSSDatabases) Add(certPath string, certificateName string) error {
	return n.executeOnPaths(func(path string) error {
		err := runCertutilCommand("-A", "-d", path, "-t", "C,,", "-n", certificateName, "-i", certPath)
		logger.Logger.Debugf("adding the certificate to the NSS database (%s): %s", path, err)

		return err
	})
}

// Delete deletes the certificate from the NSS databases.
func (n *NSSDatabases) Delete(certificateName string) error {
	return n.executeOnPaths(func(path string) error {
		err := runCertutilCommand("-D", "-d", path, "-n", certificateName)
		logger.Logger.Debugf("deleting the certificate from the NSS database (%s): %s", path, err)

		return err
	})
}

// IsKnown checks if the certificate is known by the NSS databases.
func (n *NSSDatabases) IsKnown(certificateName string) bool {
	for _, path := range n.Paths {
		err := runCertutilCommand("-V", "-d", path, "-u", "L", "-n", certificateName)
		logger.Logger.Debugf("checking if the certificate is known by the NSS database (%s): %s", path, err)

		if err != nil && !strings.Contains(err.Error(), "PR_FILE_NOT_FOUND_ERROR:") {
			logger.Logger.Errorf("failed to check if the certificate is known by the NSS databases: %v", err)
		}

		if err != nil {
			return false
		}
	}

	return true
}

// filterExistingPath filters the given paths and returns only the existing ones.
func filterExistingPath(theoricPath []string) []string {
	var dbPath []string

	for _, path := range theoricPath {
		matches, err := filepath.Glob(path)
		if err != nil {
			logger.Logger.Errorf("failed to parse nssdb pattern (%s):  %v", path, err)

			continue
		}

		// if the path is a pattern, we need to filter the dirctories not containing a database
		if strings.Contains(path, "*") {
			for _, match := range matches {
				dbFiles, _ := filepath.Glob(filepath.Join(match, databasePattern))
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

// runCertutilCommand runs the certutil command with the given arguments.
func runCertutilCommand(args ...string) error {
	bin, err := exec.LookPath("certutil")
	if err != nil {
		return fmt.Errorf("failed to find certutil binary: %w", err)
	}

	cmd := exec.Command(bin, args...)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %w", string(out), err)
	}

	return nil
}
