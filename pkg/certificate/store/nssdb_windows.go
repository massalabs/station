package store

// NSSDatabases is a struct that represents the known NSS databases on the system.
type NSSDatabases struct {
	Paths []string
}

// NewNssDatabases returns a new NSSDatabases instance.
func NewNssDatabases() (*NSSDatabases, error) {
	return nil, fmt.Errorf("not implemented")
}
