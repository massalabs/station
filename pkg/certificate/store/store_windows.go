//go:build windows
// +build windows

package store

import (
	"crypto/x509"
	"fmt"
	"syscall"
	"unsafe"
)

const (
	CRYPT_E_NOT_FOUND = 0x80092004
)

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
	if store == 0 || err != nil {
		return 0, err
	}

	return store, nil
}

func closeStore(store uintptr) error {
	if store == 0 {
		return nil
	}

	_, _, err := procCertCloseStore.Call(store, 0)
	if err != nil {
		return fmt.Errorf("failed to close windows root store: %w", err)
	}

	return nil
}

// addCertificateToStore adds the given certificate to the windows root store.
func addCertificateToStore(store uintptr, cert *x509.Certificate) error {
	if store == 0 {
		return fmt.Errorf("pointer is nil")
	}

	_, _, err := procCertAddEncodedCertificateToStore.Call(
		uintptr(store),
		uintptr(syscall.X509_ASN_ENCODING|syscall.PKCS_7_ASN_ENCODING),
		uintptr(unsafe.Pointer(&cert.Raw[0])),
		uintptr(len(cert.Raw)),
		3,
		0,
	)
	if err != nil {
		return err
	}

	return nil
}

// deleteCertificateFromStore removes the given certificate from the windows root store.
func deleteCertificateFromStore(store uintptr, cert *x509.Certificate) error {
	var certSyscall *syscall.CertContext
	var certPtr uintptr

	for {
		// Fetch next certificate
		certPtr, _, err := procCertEnumCertificatesInStore.Call(store, certPtr)
		if err != nil {
			errNumber, ok := err.(syscall.Errno)
			if ok && errNumber == CRYPT_E_NOT_FOUND {
				return ErrCertificateNotFound
			}

			return fmt.Errorf("failed to fetch certificate: %w", err)
		}

		certSyscall = (*syscall.CertContext)(unsafe.Pointer(certPtr))

		// Parse cert
		certBytes := (*[1 << 20]byte)(unsafe.Pointer(certSyscall.EncodedCert))[:certSyscall.Length]
		certX509, err := x509.ParseCertificate(certBytes)

		// Ignore parsing errors
		if err != nil || x509.SerialNumber == nil {
			continue
		}

		// Compare certificate serial numbers
		if certX509.SerialNumber.Cmp(cert.SerialNumber) == 0 {
			_, _, err = procCertDeleteCertificateFromStore.Call(certPtr)
			if err != nil {
				return fmt.Errorf("failed to delete certificate: %w", err)
			}

			return nil
		}
	}
}
