//go:build unix

package nss

import "strings"

// CertUtilService encapsulates operations on NSS database using certutil command.
type CertUtilService struct {
	runner Runner
}

// NewCertUtilService returns a new CertUtilService.
func NewCertUtilService(runner Runner) (*CertUtilService, error) {
	return &CertUtilService{runner: runner}, nil
}

// AddCA adds a certificate to the NSS database CA list.
func (s *CertUtilService) AddCA(dbPath string, certificateName string, certPath string) error {
	//nolint:wrapcheck
	return s.runner.Run("-A", "-d", dbPath, "-t", "C,,", "-n", certificateName, "-i", certPath)
}

// DeleteCA deletes a certificate from the NSS database CA list.
func (s *CertUtilService) DeleteCA(dbPath string, certificateName string) error {
	//nolint:wrapcheck
	return s.runner.Run("-D", "-d", dbPath, "-n", certificateName)
}

// IsKnownCA checks if a certificate is known by the NSS database CA list.
func (s *CertUtilService) IsKnownCA(dbPath string, certificateName string) error {
	err := s.runner.Run("-V", "-d", dbPath, "-u", "L", "-n", certificateName)
	if err != nil && !strings.Contains(err.Error(), "PR_FILE_NOT_FOUND_ERROR:") {
		//nolint:wrapcheck
		return err
	}

	return nil
}
