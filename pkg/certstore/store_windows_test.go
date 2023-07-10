//go:build windows
// +build windows

package certstore

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/sys/windows"

	"github.com/massalabs/station/pkg/certstore/mocks"
)

// loadCertificateFromFile is an helper function that loads a certificate from
// testdata/cert.pem file.
func loadCertificateFromFile() *x509.Certificate {
	// Load the certificate from the file
	certData, err := ioutil.ReadFile("testdata/cert.pem")
	if err != nil {
		panic(err)
	}

	certBlock, _ := pem.Decode(certData)
	cert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		panic(err)
	}

	return cert
}

// certInPool is an helper function that checks if a certificate is in the certificate pool.
// It compares the certificates by their raw subject.
func certInPool(cert *x509.Certificate, pool *x509.CertPool) bool {
	for _, c := range pool.Subjects() {
		if bytes.Equal(c, cert.RawSubject) {
			return true
		}
	}
	return false
}

var (
	mockCertContext = &windows.CertContext{}
	cert            = loadCertificateFromFile()
	otherError      = errors.New("create error")
)

func TestCertStore_AddCertificate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockAPI := mocks.NewMockWinAPI(mockCtrl)

	tests := []struct {
		name      string
		handler   windows.Handle
		setupMock func()
		wantErr   error
	}{
		{
			name:      "error when handler is nil",
			handler:   windows.Handle(0),
			setupMock: func() {},
			wantErr:   ErrCertStoreHandlerNotInit,
		},
		{
			name:    "error when CertCreateCertificateContext fails",
			handler: windows.Handle(1),
			setupMock: func() {
				mockAPI.EXPECT().CertCreateCertificateContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, otherError)
			},
			wantErr: otherError,
		},
		{
			name:    "error when CertAddCertificateContextToStore fails",
			handler: windows.Handle(1),
			setupMock: func() {
				mockAPI.EXPECT().CertCreateCertificateContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockCertContext, nil)
				mockAPI.EXPECT().CertAddCertificateContextToStore(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(ErrCertAlreadyExists)
			},
			wantErr: ErrCertAlreadyExists,
		},
		{
			name:    "success case",
			handler: windows.Handle(1),
			setupMock: func() {
				mockAPI.EXPECT().CertCreateCertificateContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockCertContext, nil)
				mockAPI.EXPECT().CertAddCertificateContextToStore(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &CertStore{
				handler: tt.handler,
				winAPI:  mockAPI,
			}

			tt.setupMock()

			err := store.AddCertificate(cert)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCertStore_RemoveCertificate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockAPI := mocks.NewMockWinAPI(mockCtrl)

	tests := []struct {
		name      string
		handler   windows.Handle
		setupMock func()
		wantErr   error
	}{
		{
			name:      "error when handler is nil",
			handler:   windows.Handle(0),
			setupMock: func() {},
			wantErr:   ErrCertStoreHandlerNotInit,
		},
		{
			name:    "error when converting subject to UTF16",
			handler: windows.Handle(1),
			setupMock: func() {
				mockAPI.EXPECT().UTF16PtrFromString(cert.Subject.CommonName).Return(nil, otherError)
			},
			wantErr: otherError,
		},
		{
			name:    "error when finding certificate",
			handler: windows.Handle(1),
			setupMock: func() {
				mockAPI.EXPECT().UTF16PtrFromString(cert.Subject.CommonName).Return(new(uint16), nil)
				mockAPI.EXPECT().CertFindCertificateInStore(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, ErrCertNotFound)
			},
			wantErr: ErrCertNotFound,
		},
		{
			name:    "error when deleting certificate",
			handler: windows.Handle(1),
			setupMock: func() {
				mockAPI.EXPECT().UTF16PtrFromString(cert.Subject.CommonName).Return(new(uint16), nil)
				mockAPI.EXPECT().CertFindCertificateInStore(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(mockCertContext, nil)
				mockAPI.EXPECT().CertDeleteCertificateFromStore(mockCertContext).Return(ErrUserCanceled)
			},
			wantErr: ErrUserCanceled,
		},
		{
			name:    "success case",
			handler: windows.Handle(1),
			setupMock: func() {
				mockAPI.EXPECT().UTF16PtrFromString(cert.Subject.CommonName).Return(new(uint16), nil)
				mockAPI.EXPECT().CertFindCertificateInStore(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(mockCertContext, nil)
				mockAPI.EXPECT().CertDeleteCertificateFromStore(mockCertContext).Return(nil)
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &CertStore{
				handler: tt.handler,
				winAPI:  mockAPI,
			}

			tt.setupMock()

			err := store.RemoveCertificate(cert)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCertStore_FetchCertificates(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockAPI := mocks.NewMockWinAPI(mockCtrl)

	tests := []struct {
		name       string
		handler    windows.Handle
		setupMock  func()
		wantErr    error
		poolLength int
	}{
		{
			name:      "error when handler is nil",
			handler:   windows.Handle(0),
			setupMock: func() {},
			wantErr:   ErrCertStoreHandlerNotInit,
		},
		{
			name:    "error enumerating certificates",
			handler: windows.Handle(1),
			setupMock: func() {
				mockAPI.EXPECT().CertEnumCertificatesInStore(gomock.Any(), gomock.Any()).Return(nil, otherError)
			},
			wantErr: otherError,
		},
		{
			name:    "success case - no certificate with nil",
			handler: windows.Handle(1),
			setupMock: func() {
				mockAPI.EXPECT().CertEnumCertificatesInStore(gomock.Any(), gomock.Any()).Return(nil, nil)
			},
			wantErr:    nil,
			poolLength: 0,
		},
		{
			name:    "success case - no certificate with no cert",
			handler: windows.Handle(1),
			setupMock: func() {
				mockAPI.EXPECT().CertEnumCertificatesInStore(gomock.Any(), gomock.Any()).Return(nil, ErrCertNotFound)
			},
			wantErr:    nil,
			poolLength: 0,
		},
		{
			name:    "success case",
			handler: windows.Handle(1),
			setupMock: func() {
				certBytes, err := ioutil.ReadFile("testdata/cert.pem")
				require.NoError(t, err)

				decodedBlock, _ := pem.Decode(certBytes)

				// The pointer to the byte slice
				ptr := &decodedBlock.Bytes[0]

				certContext := &windows.CertContext{
					EncodedCert: ptr,
					Length:      uint32(len(decodedBlock.Bytes)),
				}

				mockAPI.EXPECT().CertEnumCertificatesInStore(gomock.Any(), gomock.Any()).Return(certContext, nil).Times(1)
				mockAPI.EXPECT().CertEnumCertificatesInStore(gomock.Any(), gomock.Any()).Return(nil, ErrCertNotFound).Times(1)
			},
			wantErr:    nil,
			poolLength: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &CertStore{
				handler: tt.handler,
				winAPI:  mockAPI,
			}

			tt.setupMock()

			certPool, err := store.FetchCertificates()

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Len(t, certPool.Subjects(), tt.poolLength)
			}
		})
	}
}

func TestCertStore_Close(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockAPI := mocks.NewMockWinAPI(mockCtrl)

	t.Run("handler is zero", func(t *testing.T) {
		store := &CertStore{
			handler: windows.Handle(0),
			winAPI:  mockAPI,
		}

		err := store.Close(false)
		assert.NoError(t, err)
	})

	t.Run("checkNonFreedCtx is true", func(t *testing.T) {
		store := &CertStore{
			handler: windows.Handle(1),
			winAPI:  mockAPI,
		}

		mockAPI.EXPECT().CertCloseStore(windows.Handle(1), uint32(windows.CERT_CLOSE_STORE_CHECK_FLAG)).Return(nil)

		err := store.Close(true)
		assert.NoError(t, err)
	})

	t.Run("checkNonFreedCtx is false", func(t *testing.T) {
		store := &CertStore{
			handler: windows.Handle(1),
			winAPI:  mockAPI,
		}

		mockAPI.EXPECT().CertCloseStore(windows.Handle(1), uint32(0)).Return(nil)

		err := store.Close(false)
		assert.NoError(t, err)
	})
}

func TestNewCertStore(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockAPI := mocks.NewMockWinAPI(mockCtrl)

	t.Run("error when converting store name to UTF16", func(t *testing.T) {
		mockAPI.EXPECT().UTF16PtrFromString("testStore").Return(nil, errors.New("conversion error"))

		_, err := NewCertStore(mockAPI, "testStore")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unable to convert store name to UTF16")
	})

	t.Run("error when opening system certificate store", func(t *testing.T) {
		mockAPI.EXPECT().UTF16PtrFromString("testStore").Return(new(uint16), nil)
		mockAPI.EXPECT().CertOpenSystemStore(windows.Handle(0), gomock.Any()).Return(windows.Handle(0), errors.New("store opening error"))

		_, err := NewCertStore(mockAPI, "testStore")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unable to open system certificate store")
	})

	t.Run("success case", func(t *testing.T) {
		mockAPI.EXPECT().UTF16PtrFromString("testStore").Return(new(uint16), nil)
		mockAPI.EXPECT().CertOpenSystemStore(windows.Handle(0), gomock.Any()).Return(windows.Handle(1), nil)

		store, err := NewCertStore(mockAPI, "testStore")
		assert.NoError(t, err)
		assert.NotNil(t, store)
		assert.Equal(t, windows.Handle(1), store.handler)
	})
}

func TestInterpretError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected error
		errStr   string
	}{
		{
			name:     "Nil error",
			err:      nil,
			expected: nil,
			errStr:   "",
		},
		{
			name:     "CRYPT_E_NOT_FOUND error",
			err:      windows.Errno(windows.CRYPT_E_NOT_FOUND),
			expected: ErrCertNotFound,
			errStr:   "certificate not found: Cannot find object or property. (CRYPT_E_NOT_FOUND - 0x80092004)",
		},
		{
			name:     "CRYPT_E_EXISTS error",
			err:      windows.Errno(windows.CRYPT_E_EXISTS),
			expected: ErrCertAlreadyExists,
			errStr:   "certificate already exists: The object or property already exists. (CRYPT_E_EXISTS - 0x80092005)",
		},
		{
			name:     "ERROR_CANCELLED error",
			err:      windows.Errno(windows.ERROR_CANCELLED),
			expected: ErrUserCanceled,
			errStr:   "operation cancelled by user: The operation was canceled by the user.",
		},
		{
			name:     "ERROR_ACCESS_DENIED error",
			err:      windows.Errno(windows.ERROR_ACCESS_DENIED),
			expected: ErrAccessDenied,
			errStr:   "access denied: Access is denied.",
		},
		{
			name:     "Other errno error",
			err:      windows.Errno(windows.ERROR_FILE_NOT_FOUND),
			expected: windows.Errno(windows.ERROR_FILE_NOT_FOUND),
			errStr:   "The system cannot find the file specified.",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := interpretError(tc.err)
			if tc.expected != nil {
				// shall match
				// the expected error
				assert.ErrorIs(t, actual, tc.expected)
				// the initial error
				assert.ErrorAs(t, actual, &tc.err)
				//the expected error string
				assert.EqualError(t, actual, tc.errStr)
			} else {
				assert.NoError(t, actual)
			}
		})
	}
}

func TestManualCheck(t *testing.T) {
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		t.Skip("skipping test; CI environment detected")
	}

	// Initialize the certificate store.
	store, err := NewCertStore(NewWindowsImpl(), RootStore)
	assert.NoError(t, err)

	// Add the certificate to the store.
	err = store.AddCertificate(cert)
	assert.NoError(t, err)

	// Check that the added certificate is in the list.
	pool, err := store.FetchCertificates()
	assert.NoError(t, err)
	assert.True(t, certInPool(cert, pool))

	// Delete the added certificate.
	err = store.RemoveCertificate(cert)
	assert.NoError(t, err)

	// Check that the deleted certificate is no longer in the list.
	pool, err = store.FetchCertificates()
	assert.NoError(t, err)
	assert.False(t, certInPool(cert, pool))

	// Delete a non existing certificate and verify that there is an error.
	err = store.RemoveCertificate(cert)
	assert.Error(t, err)

	// Close the store checking for non freed context.
	// This is needed to avoid memory leaks.
	err = store.Close(true)
	assert.NoError(t, err)

}
