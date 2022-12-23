package cert

import "crypto/tls"

func GenerateTlsCertificate(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {

	certBytes, priv, err := GenerateCertificate(hello.ServerName)

	if err != nil {
		return nil, err
	}
	var cert tls.Certificate
	cert.Certificate = append(cert.Certificate, certBytes)
	cert.PrivateKey = priv

	return &cert, nil
}
