package http

import "testing"

func TestExtractHostname(t *testing.T) {
	tests := []struct {
		name     string
		origin   string
		expected string
	}{
		{"HTTPS URL", "https://station.massa", "station.massa"},
		{"HTTP URL", "http://localhost:3000", "localhost"},
		{"URL with port", "https://127.0.0.1:8080", "127.0.0.1"},
		{"Just hostname", "station.massa", "station.massa"},
		{"Just localhost", "localhost", "localhost"},
		{"Just IP", "127.0.0.1", "127.0.0.1"},
		{"Empty string", "", ""},
		{"Malformed URL", "not-a-url", "not-a-url"},
		{"Invalid control chars", string([]byte{0x7f, 0x80, 0x81}), ""},
		{"URL with path", "https://station.massa/path", "station.massa"},
		{"URL with query", "https://station.massa?query=value", "station.massa"},
		{"Malicious domain", "https://malicious-station.massa.com", "malicious-station.massa.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractHostname(tt.origin)
			if result != tt.expected {
				t.Errorf("extractHostname(%s) = %s, expected %s", tt.origin, result, tt.expected)
			}
		})
	}
}

func TestIsSimpleHostname(t *testing.T) {
	tests := []struct {
		name     string
		hostname string
		expected bool
	}{
		{"Valid hostname", "station.massa", true},
		{"Valid localhost", "localhost", true},
		{"Valid IP", "127.0.0.1", true},
		{"Valid with port", "localhost:3000", true},
		{"Empty string", "", false},
		{"With spaces", "station massa", false},
		{"With tabs", "station\tmassa", false},
		{"With newlines", "station\nmassa", false},
		{"With angle brackets", "station<massa", false},
		{"With quotes", "station\"massa", false},
		{"With backticks", "station`massa", false},
		{"With braces", "station{massa", false},
		{"With control chars", string([]byte{0x7f, 0x80, 0x81}), false},
		{"Starting with hyphen", "-station.massa", false},
		{"Ending with hyphen", "station.massa-", false},
		{"Starting with dot", ".station.massa", false},
		{"Ending with dot", "station.massa.", false},
		{"Consecutive dots", "station..massa", false},
		{"Valid subdomain", "api.station.massa", true},
		{"Valid with numbers", "station1.massa2", true},
		// Additional security tests for new restrictive validation
		{"With at symbol", "station@massa", false},
		{"With hash", "station#massa", false},
		{"With dollar", "station$massa", false},
		{"With percent", "station%massa", false},
		{"With ampersand", "station&massa", false},
		{"With asterisk", "station*massa", false},
		{"With plus", "station+massa", false},
		{"With equals", "station=massa", false},
		{"With question mark", "station?massa", false},
		{"With underscore", "station_massa", false},
		{"With tilde", "station~massa", false},
		{"With pipe", "station|massa", false},
		{"With backslash", "station\\massa", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsSimpleHostname(tt.hostname)
			if result != tt.expected {
				t.Errorf("isSimpleHostname(%s) = %v, expected %v", tt.hostname, result, tt.expected)
			}
		})
	}
}
