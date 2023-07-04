//go:build linux
// +build linux

package store

import (
	"crypto/x509"
	"fmt"
)

func Add(cert *x509.Certificate) error {
	return fmt.Errorf("not implemented")
}

func Delete(cert *x509.Certificate) error {
	return fmt.Errorf("not implemented")
}
