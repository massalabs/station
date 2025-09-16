package store

import (
	"crypto/x509"
	"fmt"
	"syscall"
	"unsafe"

	"github.com/massalabs/station/pkg/logger"
)

const (
	CRYPT_E_NOT_FOUND = 0x80092004

	// Certificate store addition flags
	CERT_STORE_ADD_NEW              = 1
	CERT_STORE_ADD_REPLACE_EXISTING = 3
	CERT_STORE_ADD_ALWAYS           = 4
)

// The functions below return an error even though they succeed.
// They return a pointer, if the pointer is 0, then the error is relevant.
var (
	modcrypt32                           = syscall.NewLazyDLL("crypt32.dll")
	procCertAddEncodedCertificateToStore = modcrypt32.NewProc("CertAddEncodedCertificateToStore")
	procCertCloseStore                   = modcrypt32.NewProc("CertCloseStore")
	procCertDeleteCertificateFromStore   = modcrypt32.NewProc("CertDeleteCertificateFromStore")
	procCertEnumCertificatesInStore      = modcrypt32.NewProc("CertEnumCertificatesInStore")
	procCertOpenSystemStoreW             = modcrypt32.NewProc("CertOpenSystemStoreW")
)

// Add adds the given certificate to the windows root store.
func Add(cert *x509.Certificate) error {
	logger.Debugf("Adding certificate to Windows store: Subject=%s", cert.Subject.String())

	rootStore, err := openStore()
	if err != nil {
		return fmt.Errorf("failed to open windows root store: %w", err)
	}

	err = addCertificateToStore(rootStore, cert)
	if err != nil {
		return fmt.Errorf("failed to add certificate to windows root store: %w", err)
	}

	err = closeStore(rootStore)
	if err != nil {
		return fmt.Errorf("failed to close windows root store: %w", err)
	}

	logger.Debugf("Successfully added certificate to Windows root store")
	return nil
}

// Delete deletes the given certificate from the windows root store.
func Delete(cert *x509.Certificate) error {
	rootStore, err := openStore()
	if err != nil {
		return fmt.Errorf("failed to open windows root store: %w", err)
	}

	err = deleteCertificateFromStore(rootStore, cert)
	if err != nil {
		return fmt.Errorf("failed to delete certificate from windows root store: %w", err)
	}

	err = closeStore(rootStore)
	if err != nil {
		return fmt.Errorf("failed to close windows root store: %w", err)
	}

	return nil
}

func openStore() (uintptr, error) {
	rootStr, err := syscall.UTF16PtrFromString("ROOT")
	if err != nil {
		return 0, err
	}

	store, _, err := procCertOpenSystemStoreW.Call(0, uintptr(unsafe.Pointer(rootStr)))
	if store == 0 && err != nil {
		return 0, fmt.Errorf("failed to procCertOpenSystemStoreW: %w", err)
	}

	return store, nil
}

func closeStore(store uintptr) error {
	if store == 0 {
		return nil
	}

	ret, _, err := procCertCloseStore.Call(store, 0)
	if ret == 0 && err != nil {
		return fmt.Errorf("failed to close windows root store: %v", err)
	}

	return nil
}

// addCertificateWithFlag attempts to add a certificate using the specified store flag.
func addCertificateWithFlag(store uintptr, cert *x509.Certificate, flag uintptr) (bool, error) {
	ret, _, err := procCertAddEncodedCertificateToStore.Call(
		uintptr(store),
		uintptr(syscall.X509_ASN_ENCODING|syscall.PKCS_7_ASN_ENCODING),
		uintptr(unsafe.Pointer(&cert.Raw[0])),
		uintptr(len(cert.Raw)),
		flag,
		0,
	)

	return ret != 0, err
}

// addCertificateToStore adds the given certificate to the windows root store.
func addCertificateToStore(store uintptr, cert *x509.Certificate) error {
	if store == 0 {
		return fmt.Errorf("pointer is nil")
	}

	// Try CERT_STORE_ADD_NEW first
	success, err := addCertificateWithFlag(store, cert, CERT_STORE_ADD_NEW)
	if success {
		return nil // Success
	}

	// If failed, try CERT_STORE_ADD_REPLACE_EXISTING
	logger.Debugf("Certificate addition with ADD_NEW failed, trying REPLACE_EXISTING")
	success, err = addCertificateWithFlag(store, cert, CERT_STORE_ADD_REPLACE_EXISTING)
	if success {
		logger.Debugf("Certificate successfully replaced existing")
		return nil
	}

	// Final fallback: try CERT_STORE_ADD_ALWAYS
	success, err = addCertificateWithFlag(store, cert, CERT_STORE_ADD_ALWAYS)
	if success {
		logger.Debugf("Certificate successfully added with ALWAYS flag")
		return nil
	}

	// All strategies failed
	logger.Errorf("All certificate addition strategies failed")
	if err != nil {
		return fmt.Errorf("failed adding cert: %w", err)
	}

	return fmt.Errorf("failed adding cert: all strategies returned 0")
}

// deleteCertificateFromStore removes the given certificate from the windows root store.
func deleteCertificateFromStore(store uintptr, cert *x509.Certificate) error {
	var certSyscall *syscall.CertContext
	var certPtr uintptr

	for {
		// Fetch next certificate
		certPtr, _, err := procCertEnumCertificatesInStore.Call(store, certPtr)
		if certPtr == 0 && err != nil {

			errNumber, ok := err.(syscall.Errno)
			if ok && errNumber == CRYPT_E_NOT_FOUND {
				return ErrCertificateNotFound
			}

			return fmt.Errorf("failed to enum certificates: %w", err)
		}

		certSyscall = (*syscall.CertContext)(unsafe.Pointer(certPtr))

		// Parse cert
		certBytes := (*[1 << 20]byte)(unsafe.Pointer(certSyscall.EncodedCert))[:certSyscall.Length]
		certX509, err := x509.ParseCertificate(certBytes)

		// Ignore parsing errors
		if err != nil || certX509.SerialNumber == nil {
			continue
		}

		// Compare certificate serial numbers
		if certX509.SerialNumber.Cmp(cert.SerialNumber) == 0 {
			ret, _, err := procCertDeleteCertificateFromStore.Call(certPtr)

			if ret == 0 && err != nil {
				return fmt.Errorf("failed to delete certificate: %w", err)
			}

			return nil
		}
	}
}
