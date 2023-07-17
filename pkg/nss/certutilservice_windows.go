//go:build windows
// +build windows

package nss

import "fmt"

// CertUtilService encapsulates operations on NSS database using certutil command.
type CertUtilService struct {
	runner Runner
}

// NewCertUtilService returns a new CertUtilService.
func NewCertUtilService(_ Runner) (*CertUtilService, error) {
	return nil, fmt.Errorf("not implemented")
}

// AddCA adds a certificate to the NSS database CA list.
func (c *CertUtilService) AddCA(_ string, _ string, _ string) error {
	return fmt.Errorf("not implemented")
}

// DeleteCA deletes a certificate from the NSS database CA list.
func (c *CertUtilService) DeleteCA(_ string, _ string) error {
	return fmt.Errorf("not implemented")
}

// IsKnownCA checks if a certificate is known by the NSS database CA list.
func (c *CertUtilService) IsKnownCA(_ string, _ string) error {
	return fmt.Errorf("not implemented")
}
