// Package nss provides functionality for interacting with Network Security Services (NSS) databases.
// NSS databases typically contain SSL, TLS, and other cryptographic certificates.
//
// The package includes the following key components:
//
// The CertUtilRunner struct encapsulates certutil commands. Certutil is a command-line utility
// that can create and manage certificate and key database files, and can also create many kinds of certificates. The
// Run method provides a unified way to execute any certutil command with provided arguments.
//
// The CertUtilService struct encapsulates operations on the NSS database using the
// certutil command. It uses a CertUtilRunner to execute these commands and returns an error if the command execution
// fails.
//
// The Manager struct encapsulates operations on all NSS databases within an operating system.
// The Manager struct holds a list of database paths (dbPath), an interface to the CertUtilServicer for certutil
// command execution, and a Logger interface for logging debug operations. It is used for batch operations on all NSS
// databases.
//
// These components interact to provide an interface for managing and manipulating certificates in the NSS databases.
// This includes adding, deleting, and checking for the existence of certificates across multiple database paths.
//
// Future enhancements:
//
// More functionalities can be added to each struct if required, based on the usage of the certutil tool and the
// required interactions with the NSS databases.
//
// Note:
//
// This package is currently designed to execute operations through the certutil command-line tool. In the future,
// the NSS databases could potentially be managed directly using the NSS library.
package nss
