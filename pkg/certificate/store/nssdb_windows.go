package store

import "fmt"

// NSSDatabases is a struct that represents the known NSS databases on the system.
type NSSDatabases struct {
	Paths []string
}

// NewNssDatabases returns a new NSSDatabases instance.
func NewNssDatabases() (*NSSDatabases, error) {
	return nil, fmt.Errorf("not implemented")
}

// Add adds the certificate to the NSS databases.
func (n *NSSDatabases) Add(_ string, _ string) error {
	return fmt.Errorf("not implemented")
}

// Delete deletes the certificate from the NSS databases.
func (n *NSSDatabases) Delete(_ string) error {
	return fmt.Errorf("not implemented")
}

// IsKnown checks if the certificate is known by the NSS databases.
func (n *NSSDatabases) IsKnown(_ string) bool {
	return false
}
