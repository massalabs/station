package nss

import "github.com/massalabs/station/pkg/runner"

// This file provides a CertUtilService struct that encapsulates operations on the Network Security Services (NSS)
// database using the certutil command. The NSS is a set of libraries designed to support cross-platform development
// of security-enabled client and server applications.
//
// CertUtilService uses a Runner to execute these commands and returns an error if the command execution fails.
//
// Future enhancements:
//  If more functionalities are required to be managed through the CertUtilService, additional methods can be added to
//  the CertUtilService struct.

var _ CertUtilServicer = &CertUtilService{}

// CertUtilService encapsulates operations on NSS database using certutil command.
type CertUtilService struct {
	runner runner.Runner
}

// NewCertUtilService returns a new CertUtilService.
func NewCertUtilService(runner runner.Runner) (*CertUtilService, error) {
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
	//nolint:wrapcheck
	return s.runner.Run("-V", "-d", dbPath, "-u", "L", "-n", certificateName)
}
