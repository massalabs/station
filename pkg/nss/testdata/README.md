#Test Data for NSS Database
This directory contains test data for interacting with the NSS database.

The test data were generated as follows:

1. An empty NSS database was created using certutil:

```bash
certutil -N --empty-password -d .
```

2. A new RSA private key was generated using openssl:

```bash
openssl genpkey -algorithm RSA -out privatekey.pem -pkeyopt rsa_keygen_bits:2048
```

2. A new self-signed X.509 certificate was generated from the private key:

```bash
openssl req -new -x509 -key privatekey.pem -out cert.pem -days 365 -subj "/C=FR/L=Paris/O=Massa Labs/OU=Inno team/CN=testing nss"
```

If you need to generate new test data, you can follow these same steps.