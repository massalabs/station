package api

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/massalabs/station/pkg/logger"
)

func allowedDomains() []string {
	return []string{"station.massa", "localhost", "127.0.0.1"}
}

// DomainRestrictionMiddleware checks if the request comes from an allowed domain
// for sensitive operations (like network create, delete, switch and update)
func DomainRestrictionMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" && IsRestrictedPath(r.URL.Path) {
			if !isRequestFromAllowedDomain(r) {
				logger.Warnf("Blocked operation from unauthorized domain: %s", getRequestOrigin(r))
				http.Error(w, "Forbidden: Operations restricted to authorized domains", http.StatusForbidden)
				return
			}
		}

		handler.ServeHTTP(w, r)
	})
}

func IsRestrictedPath(path string) bool {
	restrictedPaths := []string{"/network"}
	for _, restrictedPath := range restrictedPaths {
		if strings.HasPrefix(path, restrictedPath) {
			return true
		}
	}
	return false
}

func isRequestFromAllowedDomain(r *http.Request) bool {
	origin := getRequestOrigin(r)
	hostname := extractHostname(origin)

	for _, allowedDomain := range allowedDomains() {
		if hostname == allowedDomain {
			return true
		}
	}

	return false
}

func getRequestOrigin(r *http.Request) string {
	if origin := r.Header.Get("Origin"); origin != "" {
		return origin
	}

	// Check Referer header as fallback
	if referer := r.Header.Get("Referer"); referer != "" {
		return referer
	}

	// Check Host header for local requests
	if host := r.Header.Get("Host"); host != "" {
		return host
	}

	return "unknown"
}

// extractHostname safely extracts the hostname from a URL string
func extractHostname(origin string) string {
	if origin == "" {
		return ""
	}

	// Handle cases where the origin might just be a hostname without protocol
	if !strings.Contains(origin, "://") {
		if parsed, err := url.Parse("http://" + origin); err == nil {
			return parsed.Hostname()
		}

		// Fallback: if parsing fails but origin looks like a simple hostname,
		// check if it matches any allowed domain directly
		if isSimpleHostname(origin) {
			return origin
		}
		return ""
	}

	parsed, err := url.Parse(origin)
	if err != nil {
		return ""
	}

	return parsed.Hostname()
}

// isSimpleHostname checks if a string looks like a simple hostname
// (contains only valid hostname characters and no suspicious patterns)
func isSimpleHostname(s string) bool {
	if s == "" {
		return false
	}

	// Only allow valid hostname characters: letters, digits, hyphens, dots, and colons (for ports)
	for _, r := range s {
		isLetter := (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
		isDigit := r >= '0' && r <= '9'
		isAllowedPunct := r == '-' || r == '.' || r == ':'
		if !isLetter && !isDigit && !isAllowedPunct {
			return false
		}
	}

	// Must not start or end with hyphen or dot
	if strings.HasPrefix(s, "-") || strings.HasSuffix(s, "-") ||
		strings.HasPrefix(s, ".") || strings.HasSuffix(s, ".") {
		return false
	}

	// Check for consecutive dots
	if strings.Contains(s, "..") {
		return false
	}

	return true
}
