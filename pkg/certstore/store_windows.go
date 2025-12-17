//go:build windows

// This package provides a collection of operations to manage system certificate stores on Windows.
//
// Specifically, this package offers a CertStore type, which represents a system certificate store.
// This package is designed to integrate seamlessly with the native Windows API and comes with methods for creating, opening, and closing certificate stores.
// Additionally, it provides utilities for managing the certificates within these stores, such as fetching all certificates,
// adding a new certificate, and removing an existing one by its Common Name.
//
// The package leverages the x509 standard package from Go's crypto library to provide a bridge between Windows system certificate contexts and Go's x509.Certificate instances.
// It employs the golang.org/x/sys/windows package for interfacing with the underlying Windows API.
//
// For each operation, the package exposes a set of predefined errors to signify common error states,
// such as the absence of a certificate (ErrNotFound), preexistence of a certificate (ErrExists), cancellation by a user (ErrUserCanceled),
// and insufficient permissions (ErrAccessDenied).
//
// Design considerations:
// - The package is Windows-specific, as it uses the Windows API to interact with the system's certificate stores.
// - It follows best practices for exporting sentinel errors to signify specific error states.
//
// Future enhancements:
// - Prior to expanding this package's utility or moving it to a standalone repository, consider adding support for other operating systems.
// - It may be beneficial to add more features related to certificate management, depending on the broader application requirements.
package certstore

import (
	"crypto/x509"
	"errors"
	"fmt"
	"reflect"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Constants defining the system certificate store names.
const (
	RootStore = "ROOT"
	CAStore   = "CA"
	MyStore   = "MY"
	SPCStore  = "SPC"
)

// Predefined errors.
var (
	ErrCertNotFound            = errors.New("certificate not found")
	ErrCertAlreadyExists       = errors.New("certificate already exists")
	ErrUserCanceled            = errors.New("operation cancelled by user")
	ErrAccessDenied            = errors.New("access denied")
	ErrCertStoreHandlerNotInit = errors.New("certificate store handler not initialized")
)

// CertStore represents a system certificate store.
type CertStore struct {
	handler windows.Handle
	winAPI  WinAPI
}

// NewCertStore creates a new CertStore and opens the system certificate store with the specified name.
func NewCertStore(api WinAPI, name string) (*CertStore, error) {
	namePtr, err := api.UTF16PtrFromString(name)
	if err != nil {
		return nil, fmt.Errorf("unable to convert store name to UTF16: %w", err)
	}

	handler, err := api.CertOpenSystemStore(0, namePtr)
	if err != nil {
		return nil, fmt.Errorf("unable to open system certificate store: %w", err)
	}

	return &CertStore{handler: handler, winAPI: api}, nil
}

// Close closes the CertStore. If checkNonFreedCtx is true, it checks for non-freed memory.
func (s *CertStore) Close(checkNonFreedCtx bool) error {
	if s.handler == 0 {
		return nil
	}

	flags := uint32(0)
	if checkNonFreedCtx {
		flags = windows.CERT_CLOSE_STORE_CHECK_FLAG
	}

	err := s.winAPI.CertCloseStore(s.handler, flags)

	return interpretError(err)
}

// FetchCertificates returns all certificates in the store.
// The store must be initialized prior to calling this method.
func (s *CertStore) FetchCertificates() (*x509.CertPool, error) {
	if s.handler == 0 {
		return nil, ErrCertStoreHandlerNotInit
	}

	pool := x509.NewCertPool()
	var cert *windows.CertContext
	var err error

	for {
		cert, err = s.winAPI.CertEnumCertificatesInStore(s.handler, cert)
		if err != nil {
			err = interpretError(err)
			if errors.Is(err, ErrCertNotFound) {
				break
			}

			return nil, err
		}

		if cert == nil {
			break
		}

		x509Cert, err := convertCertContextToX509(cert)
		if err == nil {
			pool.AddCert(x509Cert)
		}
	}

	return pool, nil
}

// RemoveCertificate removes a certificate from the store by its Common Name.
// The store must be initialized prior to calling this method.
func (s *CertStore) RemoveCertificate(cert *x509.Certificate) error {
	if s.handler == 0 {
		return ErrCertStoreHandlerNotInit
	}

	certContextPtr, err := s.FindCertBySubject(cert.Subject.CommonName)
	if err != nil {
		return err
	}

	err = s.winAPI.CertDeleteCertificateFromStore(certContextPtr)
	if err != nil {
		return interpretError(err)
	}

	return nil
}

// AddCertificate adds a certificate to the store.
// The store must be initialized prior to calling this method.
func (s *CertStore) AddCertificate(cert *x509.Certificate) error {
	if s.handler == 0 {
		return ErrCertStoreHandlerNotInit
	}

	certContextPtr, err := s.CreateCertContext(cert)
	if err != nil {
		return err
	}

	// First try to add as new certificate
	err = s.winAPI.CertAddCertificateContextToStore(
		s.handler,
		certContextPtr,
		windows.CERT_STORE_ADD_NEW,
		nil,
	)
	if err != nil {
		interpretedErr := interpretError(err)
		// If certificate already exists (possibly expired), replace it
		if errors.Is(interpretedErr, ErrCertAlreadyExists) {
			err = s.winAPI.CertAddCertificateContextToStore(
				s.handler,
				certContextPtr,
				windows.CERT_STORE_ADD_REPLACE_EXISTING,
				nil,
			)
			if err != nil {
				return interpretError(err)
			}
		} else {
			return interpretedErr
		}
	}

	return nil
}

// interpretError translates system-specific error codes to predefined errors.
func interpretError(err error) error {
	if err == nil {
		return nil
	}

	errno, ok := err.(syscall.Errno)
	if !ok {
		return err
	}

	switch errno {
	case syscall.Errno(windows.CRYPT_E_NOT_FOUND):
		// The error code is not exported by the Windows API, so we have to format the error message.
		return fmt.Errorf("%w: %w (CRYPT_E_NOT_FOUND - 0x%X)", ErrCertNotFound, err, windows.CRYPT_E_NOT_FOUND)
	case syscall.Errno(windows.CRYPT_E_EXISTS):
		// The error code is not exported by the Windows API, so we have to format the error message.
		return fmt.Errorf("%w: %w (CRYPT_E_EXISTS - 0x%X)", ErrCertAlreadyExists, err, windows.CRYPT_E_EXISTS)
	case syscall.Errno(windows.ERROR_CANCELLED):
		return fmt.Errorf("%w: %w", ErrUserCanceled, err)
	case syscall.Errno(windows.ERROR_ACCESS_DENIED):
		return fmt.Errorf("%w: %w", ErrAccessDenied, err)
	default:
		return err
	}
}

// convertCertContextToX509 creates an x509.Certificate from a Windows cert context.
func convertCertContextToX509(ctx *windows.CertContext) (*x509.Certificate, error) {
	var der []byte

	// The byte array is manually created from the content of the cert context.
	slice := (*reflect.SliceHeader)(unsafe.Pointer(&der))
	slice.Data = uintptr(unsafe.Pointer(ctx.EncodedCert))
	slice.Len = int(ctx.Length)
	slice.Cap = int(ctx.Length)

	return x509.ParseCertificate(der)
}

// FindCertBySubject returns a certificate context by its Common Name.
func (s *CertStore) FindCertBySubject(subject string) (*windows.CertContext, error) {
	subjectPtr, err := s.winAPI.UTF16PtrFromString(subject)
	if err != nil {
		return nil, err
	}

	certContextPtr, err := s.winAPI.CertFindCertificateInStore(
		s.handler,
		windows.X509_ASN_ENCODING|windows.PKCS_7_ASN_ENCODING,
		0,
		windows.CERT_FIND_SUBJECT_STR_W,
		unsafe.Pointer(subjectPtr),
		nil,
	)
	if err != nil {
		return nil, interpretError(err)
	}

	return certContextPtr, nil
}

// CreateCertContext creates a new system certificate context from an x509.Certificate.
func (s *CertStore) CreateCertContext(cert *x509.Certificate) (*windows.CertContext, error) {
	certContextPtr, err := s.winAPI.CertCreateCertificateContext(
		windows.X509_ASN_ENCODING|windows.PKCS_7_ASN_ENCODING,
		&cert.Raw[0],
		uint32(len(cert.Raw)),
	)
	if err != nil {
		return nil, interpretError(err)
	}

	return certContextPtr, nil
}
