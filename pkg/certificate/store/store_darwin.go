//go:build darwin
// +build darwin

package store

import (
	"crypto/x509"
	"fmt"
)

func AddToOS(_ *x509.Certificate) error {
	return fmt.Errorf("not implemented")
}

func DeleteFromOS(_ *x509.Certificate) error {
	return fmt.Errorf("not implemented")
}
