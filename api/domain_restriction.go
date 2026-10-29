package api

import (
	"net/http"
	"strings"

	stationHttp "github.com/massalabs/station/pkg/http"
	"github.com/massalabs/station/pkg/logger"
)

func allowedDomains() []string {
	return []string{"station.massa", "localhost", "127.0.0.1"}
}

// DomainRestrictionMiddleware checks if the request comes from an allowed domain
// for sensitive operations (like network create, delete, switch and update)
// and for sensitive data (like logs), which CORS would otherwise expose to any website.
func DomainRestrictionMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if (r.Method != "GET" && IsRestrictedPath(r.URL.Path)) || IsRestrictedReadPath(r.URL.Path) {
			if !isRequestFromAllowedDomain(r) {
				logger.Warnf("Blocked operation from unauthorized domain: %s", stationHttp.GetRequestOrigin(r))
				http.Error(w, "Forbidden: Operations restricted to authorized domains", http.StatusForbidden)
				return
			}
		}

		handler.ServeHTTP(w, r)
	})
}

func IsRestrictedPath(path string) bool {
	restrictedPaths := []string{"/network", "/plugin-manager"}
	for _, restrictedPath := range restrictedPaths {
		if strings.HasPrefix(path, restrictedPath) {
			return true
		}
	}
	return false
}

// IsRestrictedReadPath returns true for paths that must be restricted to allowed domains for all methods.
func IsRestrictedReadPath(path string) bool {
	return strings.TrimSuffix(path, "/") == "/logs"
}

func isRequestFromAllowedDomain(r *http.Request) bool {
	origin := stationHttp.GetRequestOrigin(r)
	hostname := stationHttp.ExtractHostname(origin)

	for _, allowedDomain := range allowedDomains() {
		if hostname == allowedDomain {
			return true
		}
	}

	return false
}
