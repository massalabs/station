package nss

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCertUtilServicer is a mock of CertUtilServicer.
type MockCertUtilServicer struct {
	mock.Mock
}

func (m *MockCertUtilServicer) AddCA(path string, certName string, certPath string) error {
	args := m.Called(path, certName, certPath)

	//nolint:wrapcheck
	return args.Error(0)
}

func (m *MockCertUtilServicer) DeleteCA(path string, certName string) error {
	args := m.Called(path, certName)

	//nolint:wrapcheck
	return args.Error(0)
}

func (m *MockCertUtilServicer) IsKnownCA(path string, certName string) error {
	args := m.Called(path, certName)

	//nolint:wrapcheck
	return args.Error(0)
}

// MockLogger is a mock of Logger.
type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Debugf(msg string, args ...interface{}) {
	m.Called(msg, args)
}

func (m *MockLogger) Errorf(msg string, args ...interface{}) {
	m.Called(msg, args)
}

func TestManager_AddCA(t *testing.T) {
	mockCertUtilServicer := new(MockCertUtilServicer)
	mockLogger := new(MockLogger)
	m := NewManager([]string{"/path/to/db"}, mockCertUtilServicer, mockLogger)

	mockCertUtilServicer.On("AddCA", "/path/to/db", "testCert", "/path/to/cert").Return(nil)
	mockLogger.On("Debugf", mock.Anything, mock.Anything)

	err := m.AddCA("testCert", "/path/to/cert")

	assert.NoError(t, err)
	mockCertUtilServicer.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestManager_DeleteCA(t *testing.T) {
	mockCertUtilServicer := new(MockCertUtilServicer)
	mockLogger := new(MockLogger)
	m := NewManager([]string{"/path/to/db"}, mockCertUtilServicer, mockLogger)

	mockCertUtilServicer.On("DeleteCA", "/path/to/db", "testCert").Return(nil)
	mockLogger.On("Debugf", mock.Anything, mock.Anything)

	err := m.DeleteCA("testCert")

	assert.NoError(t, err)
	mockCertUtilServicer.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestManager_HasCA(t *testing.T) {
	mockCertUtilServicer := new(MockCertUtilServicer)
	mockLogger := new(MockLogger)
	m := NewManager([]string{"/path/to/db"}, mockCertUtilServicer, mockLogger)

	mockCertUtilServicer.On("IsKnownCA", "/path/to/db", "testCert").Return(nil)
	mockLogger.On("Debugf", mock.Anything, mock.Anything)

	result := m.HasCA("testCert")

	assert.True(t, result)
	mockCertUtilServicer.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestManager_AddCA_Error(t *testing.T) {
	mockCertUtilServicer := new(MockCertUtilServicer)
	mockLogger := new(MockLogger)
	m := NewManager([]string{"/path/to/db"}, mockCertUtilServicer, mockLogger)

	mockCertUtilServicer.On("AddCA", "/path/to/db", "testCert", "/path/to/cert").Return(errors.New("mock error"))
	mockLogger.On("Debugf", mock.Anything, mock.Anything)

	err := m.AddCA("testCert", "/path/to/cert")

	assert.Error(t, err)
	mockCertUtilServicer.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestManager_DeleteCA_Error(t *testing.T) {
	mockCertUtilServicer := new(MockCertUtilServicer)
	mockLogger := new(MockLogger)
	m := NewManager([]string{"/path/to/db"}, mockCertUtilServicer, mockLogger)

	mockCertUtilServicer.On("DeleteCA", "/path/to/db", "testCert").Return(errors.New("mock error"))
	mockLogger.On("Debugf", mock.Anything, mock.Anything)

	err := m.DeleteCA("testCert")

	assert.Error(t, err)
	mockCertUtilServicer.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestManager_HasCA_Error(t *testing.T) {
	mockCertUtilServicer := new(MockCertUtilServicer)
	mockLogger := new(MockLogger)
	m := NewManager([]string{"/path/to/db"}, mockCertUtilServicer, mockLogger)

	mockCertUtilServicer.On("IsKnownCA", "/path/to/db", "testCert").Return(errors.New("mock error"))
	mockLogger.On("Debugf", mock.Anything, mock.Anything)

	result := m.HasCA("testCert")

	assert.False(t, result)
	mockCertUtilServicer.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestExpandAndFilterPaths(t *testing.T) {
	var (
		pathTo                = filepath.Join("path", "to")
		pathToPattern         = filepath.Join("path", "to", "*")
		pathToDir1            = filepath.Join("path", "to", "dir1")
		pathToDir1Cert9       = filepath.Join("path", "to", "dir1", "cert9.pem")
		pathToDir1CertPattern = filepath.Join("path", "to", "dir1", databasePattern)
		pathToDir2            = filepath.Join("path", "to", "dir2")
	)

	tests := []struct {
		name     string
		input    []string
		mockGlob func(pattern string) ([]string, error)
		want     []string
	}{
		{
			name:  "normal operation",
			input: []string{pathTo, pathToDir1, pathToDir2},
			mockGlob: func(pattern string) ([]string, error) {
				if pattern == pathToDir1 {
					return []string{pathToDir1}, nil
				}
				if pattern == pathToDir2 {
					return []string{}, nil
				}

				return nil, nil
			},
			want: []string{pathToDir1},
		},
		{
			name:  "dynamic path operation",
			input: []string{pathToPattern},
			mockGlob: func(pattern string) ([]string, error) {
				if pattern == pathToPattern {
					return []string{pathToDir1, pathToDir2}, nil
				}
				if pattern == pathToDir1CertPattern {
					return []string{pathToDir1Cert9}, nil
				}
				if pattern == pathToDir1CertPattern {
					return []string{}, nil
				}

				return nil, nil
			},
			want: []string{pathToDir1},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			manager := NewManager([]string{}, nil, nil)
			got := manager.expandAndFilterPaths(tc.mockGlob, tc.input)
			assert.Equal(t, tc.want, got)
		})
	}
}
