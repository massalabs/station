package certificate

import (
	"testing"
)

const (
	validCert = `-----BEGIN CERTIFICATE-----
    dG90bwo=
-----END CERTIFICATE-----`
	validPrivateKey = `-----BEGIN PRIVATE KEY-----
	dG90bwo=
-----END PRIVATE KEY-----`
	invalidPEMType = `-----BEGIN INVALID-----
	dG90bwo=
-----END INVALID-----`
	invalidBase64 = `-----BEGIN INVALID-----
	%
-----END INVALID-----`
)

func TestDecodePEM(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		wantErr bool
	}{
		{
			name:    "Valid certificate",
			data:    validCert,
			wantErr: false,
		},
		{
			name:    "Valid private key",
			data:    validPrivateKey,
			wantErr: false,
		},
		{
			name:    "Invalid data",
			data:    invalidBase64,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DecodePEM([]byte(tt.data))
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodePEM() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDecodeExpectedPEM(t *testing.T) {
	tests := []struct {
		name     string
		data     string
		expected PemType
		wantErr  bool
	}{
		{
			name:     "Valid certificate",
			data:     validCert,
			expected: Certificate,
			wantErr:  false,
		},
		{
			name:     "Valid private key",
			data:     validPrivateKey,
			expected: PrivateKey,
			wantErr:  false,
		},
		{
			name:     "Inconsistent PEM type",
			data:     validCert,
			expected: PrivateKey,
			wantErr:  true,
		},
		{
			name:     "Inconsistent PEM type",
			data:     invalidPEMType,
			expected: PrivateKey,
			wantErr:  true,
		},
		{
			name:     "Invalid data",
			data:     invalidBase64,
			expected: Certificate,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DecodeExpectedPEM([]byte(tt.data), tt.expected)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeWantedPEM() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPemTypeIsValid(t *testing.T) {
	tests := []struct {
		name string
		p    PemType
		want bool
	}{
		{
			name: "Valid Certificate",
			p:    Certificate,
			want: true,
		},
		{
			name: "Valid PrivateKey",
			p:    PrivateKey,
			want: true,
		},
		{
			name: "Valid Certificate Request",
			p:    CertificateRequest,
			want: true,
		},
		{
			name: "Valid X509 CRL",
			p:    X509CRL,
			want: true,
		},
		{
			name: "Invalid PemType (too big)",
			p:    PrivateKey + 1,
			want: false,
		},
		{
			name: "Invalid PemType (too small)",
			p:    0,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.IsValid(); got != tt.want {
				t.Errorf("PemType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPemTypeString(t *testing.T) {
	tests := []struct {
		name string
		p    PemType
		want string
	}{
		{
			name: "Certificate String",
			p:    Certificate,
			want: "CERTIFICATE",
		},
		{
			name: "PrivateKey String",
			p:    PrivateKey,
			want: "PRIVATE KEY",
		},
		{
			name: "Certificate Request String",
			p:    CertificateRequest,
			want: "CERTIFICATE REQUEST",
		},
		{
			name: "X509 CRL String",
			p:    X509CRL,
			want: "X509 CRL",
		},
		{
			name: "Invalid PemType (too big)",
			p:    PrivateKey + 1,
			want: "",
		},
		{
			name: "Invalid PemType (too small)",
			p:    0,
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.String(); got != tt.want {
				t.Errorf("PemType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
