package nss

import (
	"errors"
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
	tests := []struct {
		name      string
		input     []string
		mockGlob  func(pattern string) ([]string, error)
		want      []string
		expectErr bool
	}{
		{
			name:  "normal operation",
			input: []string{"/path/to/", "/path/to/dir1/", "/path/to/dir2/"},
			mockGlob: func(pattern string) ([]string, error) {
				if pattern == "/path/to/dir1/" {
					return []string{"/path/to/dir1/"}, nil
				}
				if pattern == "/path/to/dir2/" {
					return []string{}, nil
				}
				return nil, nil
			},
			want: []string{"/path/to/dir1/"},
		},
		{
			name:  "dynamic path operation",
			input: []string{"/path/to/*"},
			mockGlob: func(pattern string) ([]string, error) {
				if pattern == "/path/to/*" {
					return []string{"/path/to/dir1", "/path/to/dir2"}, nil
				}
				if pattern == "/path/to/dir1/cert*.db" {
					return []string{"/path/to/dir1/cert9.db"}, nil
				}
				if pattern == "/path/to/dir2/cert*.db" {
					return []string{}, nil
				}
				return nil, nil
			},
			want: []string{"/path/to/dir1"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			manager := NewManager([]string{}, nil, nil)
			got, err := manager.expandAndFilterPaths(tc.mockGlob, tc.input)
			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.want, got)
		})
	}
}
