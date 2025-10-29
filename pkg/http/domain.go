package http

import (
	"net/http"
	"net/url"
	"strings"
)

func GetRequestOrigin(r *http.Request) string {
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
func ExtractHostname(origin string) string {
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
		if IsSimpleHostname(origin) {
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
func IsSimpleHostname(s string) bool {
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
