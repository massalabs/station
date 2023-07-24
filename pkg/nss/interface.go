package nss

// Logger is the interface used by the NSS package to log messages.
type Logger interface {
	Debugf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
}

// CertUtilServicer encapsulates operations on NSS database using certutil command.
type CertUtilServicer interface {
	// AddCA adds a certificate to the NSS database CA list.
	AddCA(dbPath string, certificateName string, certPath string) error
	// DeleteCA deletes a certificate from the NSS database CA list.
	DeleteCA(dbPath string, certificateName string) error
	// IsKnownCA checks if a certificate is known by the NSS database CA list.
	IsKnownCA(dbPath string, certificateName string) error
}

// Runner is the interface used by the NSS package to run commands.
type Runner interface {
	// Run runs the given command and returns the combined output of stdout and stderr.
	Run(args ...string) error
}
