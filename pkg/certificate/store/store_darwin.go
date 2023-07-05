//go:build darwin
// +build darwin

package store

import (
	"crypto/x509"
	"fmt"
)

func Add(_ *x509.Certificate) error {
	return fmt.Errorf("not implemented")
}

func Delete(_ *x509.Certificate) error {
	return fmt.Errorf("not implemented")
}
