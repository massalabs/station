package dnshelper

import (
	"fmt"
	"html/template"

	"github.com/massalabs/station/pkg/convert"
	"github.com/massalabs/station/pkg/node/base58"
)

const (
	// Indexes of data in website name key.
	indexOfWebsiteAddress     = 0 // Index of website Address in the dnsValue array
	indexOfWebsiteDescription = 2 // Index of website Description in the dnsValue array
)

// AddressAndDescription fetches the website address and its description from the DNS entry.
func AddressAndDescription(dnsValue []byte) (string, string, error) {
	dnsRecords := convert.ByteToStringArray(dnsValue)
	if len(dnsRecords) <= indexOfWebsiteAddress {
		return "", "", fmt.Errorf("invalid website value: missing website address")
	}

	address := dnsRecords[indexOfWebsiteAddress]
	if !IsValidAddress(address) {
		return "", "", fmt.Errorf("invalid website address: %s", address)
	}

	description := ""

	if len(dnsRecords) > indexOfWebsiteDescription {
		unsafeDescription := dnsRecords[indexOfWebsiteDescription]

		// Prevent XSS by escaping special characters in websiteDescription
		// see
		// https://cheatsheetseries.owasp.org/cheatsheets/Cross_Site_Scripting_
		// Prevention_Cheat_Sheet.html#output-encoding-for-html-contexts
		description = template.HTMLEscapeString(unsafeDescription)
	}

	return address, description, nil
}

// IsValidAddress checks if the address is valid based on the prefix rule, non-empty rule, and successful decoding.
func IsValidAddress(addr string) bool {
	if addr == "" {
		return false
	}

	addressPrefix := addr[:2]
	addressWithoutPrefix := addr[2:]

	if addressPrefix == "AS" && len(addressWithoutPrefix) > 0 {
		_, _, err := base58.VersionedCheckDecode(addressWithoutPrefix)
		if err == nil {
			return true
		}
	}

	return false
}
