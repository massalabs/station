package api

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/massalabs/station/pkg/logger"
)

func TestMain(m *testing.M) {
	// Initialize logger for tests
	tempDir, err := os.MkdirTemp("", "test-logs")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tempDir) // nolint: errcheck

	logPath := filepath.Join(tempDir, "test.log")
	if err := logger.InitializeGlobal(logPath); err != nil {
		panic(err)
	}

	// Run tests
	os.Exit(m.Run())
}

func TestDomainRestrictionMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		origin         string
		referer        string
		host           string
		expectedStatus int
		shouldBlock    bool
	}{
		{
			name:           "Allowed domain - station.massa",
			method:         "POST",
			path:           "/network/create",
			origin:         "https://station.massa",
			expectedStatus: http.StatusOK,
			shouldBlock:    false,
		},
		{
			name:           "Allowed domain - localhost",
			method:         "DELETE",
			path:           "/network/delete/testnet",
			origin:         "http://localhost:3000",
			expectedStatus: http.StatusOK,
			shouldBlock:    false,
		},
		{
			name:           "Allowed domain - 127.0.0.1",
			method:         "PUT",
			path:           "/network/switch/mainnet",
			origin:         "http://127.0.0.1:8080",
			expectedStatus: http.StatusOK,
			shouldBlock:    false,
		},
		{
			name:           "Blocked domain - external site",
			method:         "POST",
			path:           "/network/create",
			origin:         "https://malicious-site.com",
			expectedStatus: http.StatusForbidden,
			shouldBlock:    true,
		},
		{
			name:           "Blocked domain - no origin header",
			method:         "DELETE",
			path:           "/network/delete/testnet",
			origin:         "",
			expectedStatus: http.StatusForbidden,
			shouldBlock:    true,
		},
		{
			name:           "Non-network endpoint - should pass through",
			method:         "GET",
			path:           "/api/version",
			origin:         "https://malicious-site.com",
			expectedStatus: http.StatusOK,
			shouldBlock:    false,
		},
		{
			name:           "Fallback to referer header",
			method:         "PUT",
			path:           "/network/switch/testnet",
			origin:         "",
			referer:        "https://station.massa/network",
			expectedStatus: http.StatusOK,
			shouldBlock:    false,
		},
		{
			name:           "Fallback to host header",
			method:         "POST",
			path:           "/network/create",
			origin:         "",
			referer:        "",
			host:           "localhost:3000",
			expectedStatus: http.StatusOK,
			shouldBlock:    false,
		},
		{
			name:           "GET /network from any domain - should pass through",
			method:         "GET",
			path:           "/network",
			origin:         "https://malicious-site.com",
			expectedStatus: http.StatusOK,
			shouldBlock:    false,
		},
		{
			name:           "GET /network with no headers - should pass through",
			method:         "GET",
			path:           "/network",
			origin:         "",
			referer:        "",
			host:           "",
			expectedStatus: http.StatusOK,
			shouldBlock:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test handler that returns 200 OK
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			// Apply the domain restriction middleware
			middleware := DomainRestrictionMiddleware(testHandler)

			// Create a test request
			req := httptest.NewRequest(tt.method, tt.path, nil)
			if tt.origin != "" {
				req.Header.Set("Origin", tt.origin)
			}
			if tt.referer != "" {
				req.Header.Set("Referer", tt.referer)
			}
			if tt.host != "" {
				req.Header.Set("Host", tt.host)
			}

			// Create a response recorder
			rr := httptest.NewRecorder()

			// Serve the request
			middleware.ServeHTTP(rr, req)

			// Check the response
			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			// Check if the request was blocked
			if tt.shouldBlock && rr.Code != http.StatusForbidden {
				t.Errorf("expected request to be blocked with 403, got %d", rr.Code)
			}

			if !tt.shouldBlock && rr.Code != http.StatusOK {
				t.Errorf("expected request to pass through with 200, got %d", rr.Code)
			}
		})
	}
}

func TestIsRestrictedPath(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{"/network", true},
		{"/network/create", true},
		{"/network/delete/testnet", true},
		{"/network/switch/mainnet", true},
		{"/api/version", false},
		{"/plugin/install", false},
		{"/", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := IsRestrictedPath(tt.path)
			if result != tt.expected {
				t.Errorf("isRestrictedPath(%s) = %v, expected %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestIsRequestFromAllowedDomain(t *testing.T) {
	tests := []struct {
		name     string
		origin   string
		expected bool
	}{
		{"station.massa", "https://station.massa", true},
		{"localhost", "http://localhost:3000", true},
		{"127.0.0.1", "http://127.0.0.1:8080", true},
		{"malicious site", "https://malicious-site.com", false},
		{"empty origin", "", false},
		{"subdomain should be blocked", "https://app.station.massa", false},
		{"port number", "http://localhost:1234", true},
		// Security tests - these should all be blocked with the new implementation
		{"malicious domain containing allowed domain", "https://malicious-station.massa.com", false},
		{"malicious domain with allowed domain as prefix", "https://station.massa.evil.com", false},
		{"malicious domain with allowed domain in path", "https://evil.com/station.massa", false},
		{"malicious localhost lookalike", "https://localhost.evil.com", false},
		{"IP address lookalike", "https://127.0.0.1.evil.com", false},
		// Edge cases
		{"just hostname without protocol", "station.massa", true},
		{"just localhost without protocol", "localhost", true},
		{"just IP without protocol", "127.0.0.1", true},
		{"malformed URL", "not-a-url", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			if tt.origin != "" {
				req.Header.Set("Origin", tt.origin)
			}

			result := isRequestFromAllowedDomain(req)
			if result != tt.expected {
				t.Errorf("isRequestFromAllowedDomain() = %v, expected %v for origin %s", result, tt.expected, tt.origin)
			}
		})
	}
}
