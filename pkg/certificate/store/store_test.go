package store

import (
	"crypto/x509"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCAUniqueFilename(t *testing.T) {
	filename := CAUniqueFilename(&x509.Certificate{SerialNumber: big.NewInt(1)})
	assert.Equal(t, "MassaLabs_CA_1.cert", filename)
}
