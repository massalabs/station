package certificate

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"fmt"
)

const privateKeySizeInBits = 2048

func GenerateTLS(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	priv, err := rsa.GenerateKey(rand.Reader, privateKeySizeInBits)
	if err != nil {
		// here we panic because we have no way to get the error in the terminal...
		panic(fmt.Errorf("while loading CA certificate: %w", err))
	}

	certBytes, priv, err := GenerateSignedCertificate(hello.ServerName, priv)
	if err != nil {
		panic(err) // here we panic because we have no way to get the error in the terminal...
	}

	var cert tls.Certificate

	cert.Certificate = append(cert.Certificate, certBytes)
	cert.PrivateKey = priv

	return &cert, nil
}
