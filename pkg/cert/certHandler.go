package cert

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

const rootName = "rootCA.pem"
const rootKeyName = "rootCA-key.pem"

// CAROOT ...
var CAROOT = getCAROOT()

// CA certificate and key.
var caCert *x509.Certificate
var caKey crypto.PrivateKey

// Private Key for the certificate.
var priv *rsa.PrivateKey

// Static data for the certificate.
var tempCertificate = x509.Certificate{
	SerialNumber: randomSerialNumber(),
	Subject: pkix.Name{
		Organization: []string{"thyra dynamically generated"},
	},
	ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
	KeyUsage:    x509.KeyUsageDigitalSignature,
}

func getCa() (*x509.Certificate, crypto.PrivateKey) {
	// Get the CA root certificate.
	if CAROOT == "" {
		log.Fatalln("ERROR: failed to find the default CA location, set one as the CAROOT env var")
	}
	// Load the CA certificate and key.
	return loadCA()
}

func generatePrivAndPubKey() (*rsa.PrivateKey, crypto.PublicKey) {
	// Generate a private key.
	var err error

	if priv != nil {
		return priv, priv.Public()
	}

	const random = 2048

	priv, err = rsa.GenerateKey(rand.Reader, random)

	fatalIfErr(err, "failed to generate private key")

	return priv, priv.Public()
}

func GenerateCertificate(serverName string) ([]byte, crypto.PrivateKey, error) {
	caCert, caKey = getCa()

	if len(serverName) == 0 {
		panic("error:ServerName is Empty.")
	}

	priv, pub := generatePrivAndPubKey()

	// Create the template for the certificate.
	expiration := time.Now().AddDate(0, 1, 0)
	tpl := &tempCertificate
	tpl = &x509.Certificate{
		Subject: pkix.Name{
			CommonName: serverName,
		},
		NotBefore: time.Now(),
		NotAfter:  expiration,
		DNSNames:  append(tpl.DNSNames, serverName),
	}
	// Create the certificate.
	cert, err := x509.CreateCertificate(rand.Reader, tpl, caCert, pub, caKey)

	return cert, priv, err
}

// Get the CA root certificate depending on your OS.
func getCAROOT() string {
	if env := os.Getenv("CAROOT"); env != "" {
		return env
	}

	var dir string

	switch {
	case runtime.GOOS == "windows":
		dir = os.Getenv("LocalAppData")

	case os.Getenv("XDG_DATA_HOME") != "":
		dir = os.Getenv("XDG_DATA_HOME")

	case runtime.GOOS == "darwin":
		dir = os.Getenv("HOME")
		if dir == "" {
			// $HOME is not set.
			// return empty string when CAROOT is not set.
			return ""
		}

		dir = filepath.Join(dir, "Library", "Application Support")

	case runtime.GOOS == "linux":
		dir = os.Getenv("HOME")
		if dir == "" {
			// $HOME is not set.
			// return empty string when CAROOT is not set.
			return ""
		}
	default: // Unix.
		log.Fatalln("ERROR: OS not supported please contact support")
	}

	return filepath.Join(dir, "mkcert")
}

// Get the CA root key or create a new one.
func loadCA() (*x509.Certificate, crypto.PrivateKey) {
	if !pathExists(filepath.Join(CAROOT, rootName)) {
		log.Fatalln("ERROR: failed to find the CA Root path at: " + filepath.Join(CAROOT, rootName) + "")
	}

	certPEMBlock, err := os.ReadFile(filepath.Join(CAROOT, rootName))

	fatalIfErr(err, "failed to read the CA certificate")

	certDERBlock, _ := pem.Decode(certPEMBlock)
	if certDERBlock == nil || certDERBlock.Type != "CERTIFICATE" {
		log.Fatalln("ERROR: failed to read the CA certificate: unexpected content")
	}

	caCert, err := x509.ParseCertificate(certDERBlock.Bytes)

	fatalIfErr(err, "failed to parse the CA certificate")

	if !pathExists(filepath.Join(CAROOT, rootKeyName)) {
		log.Fatalln("ERROR: failed to find the CA RootKeyName path at: " + filepath.Join(CAROOT, rootKeyName) + "")
	}

	keyPEMBlock, err := ioutil.ReadFile(filepath.Join(CAROOT, rootKeyName))

	fatalIfErr(err, "failed to read the CA key")

	keyDERBlock, _ := pem.Decode(keyPEMBlock)
	if keyDERBlock == nil || keyDERBlock.Type != "PRIVATE KEY" {
		log.Fatalln("ERROR: failed to read the CA key: unexpected content")
	}
	caKey, err = x509.ParsePKCS8PrivateKey(keyDERBlock.Bytes)

	fatalIfErr(err, "failed to parse the CA key")

	return caCert, caKey
}

// randomSerialNumber generates a random serial number for the certificate.
func randomSerialNumber() *big.Int {
	const bits = 128 // 128 bits
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), bits)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	fatalIfErr(err, "failed to generate serial number")

	return serialNumber
}

// return true if path exist.
func pathExists(path string) bool {
	_, err := os.Stat(path)

	return err == nil
}

func fatalIfErr(err error, msg string) {
	if err != nil {
		log.Fatalf("ERROR: %s: %s", msg, err)
	}
}
