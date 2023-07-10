//go:build windows
// +build windows

package certstore

// This file provides the WinAPI interface and an its implementation for Windows using golang.org/x/sys/windows package.

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

// WinAPI is an interface for Windows API certificate functions.
// This interface is used for mocking in tests.
type WinAPI interface {
	UTF16PtrFromString(s string) (*uint16, error)
	CertOpenSystemStore(handle windows.Handle, name *uint16) (windows.Handle, error)
	CertCloseStore(handle windows.Handle, flags uint32) error
	CertEnumCertificatesInStore(store windows.Handle, prevContext *windows.CertContext) (*windows.CertContext, error)
	CertFindCertificateInStore(store windows.Handle, certEncodingType uint32, findFlags uint32, findType uint32, findPara unsafe.Pointer, prevCertContext *windows.CertContext) (*windows.CertContext, error)
	CertDeleteCertificateFromStore(certContext *windows.CertContext) error
	CertCreateCertificateContext(certEncodingType uint32, certEncoded *byte, certEncodedLen uint32) (*windows.CertContext, error)
	CertAddCertificateContextToStore(store windows.Handle, certContext *windows.CertContext, addDisposition uint32, storeContext **windows.CertContext) error
}

// WindowsImpl is an implementation of WinAPI interface.
// It uses golang.org/x/sys/windows package.
type WindowsImpl struct{}

func (WindowsImpl) UTF16PtrFromString(s string) (*uint16, error) {
	return windows.UTF16PtrFromString(s)
}

func (WindowsImpl) CertOpenSystemStore(handle windows.Handle, name *uint16) (windows.Handle, error) {
	return windows.CertOpenSystemStore(handle, name)
}

func (WindowsImpl) CertCloseStore(handle windows.Handle, flags uint32) error {
	return windows.CertCloseStore(handle, flags)
}

func (WindowsImpl) CertEnumCertificatesInStore(store windows.Handle, prevContext *windows.CertContext) (*windows.CertContext, error) {
	return windows.CertEnumCertificatesInStore(store, prevContext)
}

func (WindowsImpl) CertFindCertificateInStore(store windows.Handle, certEncodingType uint32, findFlags uint32, findType uint32, findPara unsafe.Pointer, prevCertContext *windows.CertContext) (*windows.CertContext, error) {
	return windows.CertFindCertificateInStore(store, certEncodingType, findFlags, findType, findPara, prevCertContext)
}

func (WindowsImpl) CertDeleteCertificateFromStore(certContext *windows.CertContext) error {
	return windows.CertDeleteCertificateFromStore(certContext)
}

func (WindowsImpl) CertCreateCertificateContext(certEncodingType uint32, certEncoded *byte, certEncodedLen uint32) (*windows.CertContext, error) {
	return windows.CertCreateCertificateContext(certEncodingType, certEncoded, certEncodedLen)
}

func (WindowsImpl) CertAddCertificateContextToStore(store windows.Handle, certContext *windows.CertContext, addDisposition uint32, storeContext **windows.CertContext) error {
	return windows.CertAddCertificateContextToStore(store, certContext, addDisposition, storeContext)
}

var _ WinAPI = WindowsImpl{}

// NewWindowsImpl creates a new WindowsImpl.
func NewWindowsImpl() *WindowsImpl {
	return &WindowsImpl{}
}
