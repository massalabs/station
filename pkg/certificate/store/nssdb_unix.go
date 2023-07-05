//go:build unix

package store

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/massalabs/station/pkg/logger"
)

var nssDBPaths = []string{
	"/etc/pki/nssdb/",
	os.Getenv("HOME") + "/.pki/nssdb/",
	os.Getenv("HOME") + "/snap/chromium/current/.pki/nssdb/",
	os.Getenv("HOME") + "/.mozilla/firefox/*",
	os.Getenv("HOME") + "/snap/firefox/common/.mozilla/firefox/*",
}

const databasePattern = "cert*.db"

// NSSDatabases is a struct that represents the known NSS databases on the system.
type NSSDatabases struct {
	Paths []string
}

// NewNssDatabases returns a new NSSDatabases instance.
func NewNssDatabases() *NSSDatabases {
	return &NSSDatabases{
		Paths: findDB(),
	}
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
		return runCertutilCommand("-A", "-d", path, "-t", "C,,", "-n", certificateName, "-i", certPath)
	})
}

// Delete deletes the certificate from the NSS databases.
func (n *NSSDatabases) Delete(certificateName string) error {
	return n.executeOnPaths(func(path string) error {
		return runCertutilCommand("-D", "-d", path, "-n", certificateName)
	})
}

// IsKnown checks if the certificate is known by the NSS databases.
func (n *NSSDatabases) IsKnown(certificateName string) bool {
	err := n.executeOnPaths(func(path string) error {
		return runCertutilCommand("-V", "-d", path, "-u", "L", "-n", certificateName)
	})

	return err == nil
}

// List lists the recognized NSS databases.
func findDB() []string {
	var dbPath []string

	for _, path := range nssDBPaths {
		matches, err := filepath.Glob(path)
		if err != nil {
			log.Fatalf("failed to parse nssdb pattern (%s):  %v", path, err)

			continue
		}

		//if the path is a pattern, we need to filter the dirctories not containing a database
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
		return err
	}

	cmd := exec.Command(bin, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Logger.Infoln("failed to run certutil command (%s): %v", cmd, err)
		logger.Logger.Debugln("certutil output: %s", out)

		return err
	}

	return nil
}
