package dnshelper

import (
	"fmt"
	"html/template"

	"github.com/massalabs/station/pkg/convert"
	"github.com/massalabs/station/pkg/node"
)

const (
	metaKey = "META"

	// Indexes of data in website name key.
	indexOfWebsiteAddress     = 0 // Index of website Address in the dnsValue array
	indexOfWebsiteDescription = 2 // Index of website Description in the dnsValue array
)

// AddressAndDescription fetch the website address and it's description from the DNS entry.
func AddressAndDescription(dnsValue []byte) (string, string, error) {
	// In dnsRecords we have 3 values respecting the following order:
	// websiteAddress, ownerAddress, and finally websitedescription
	// here we retrieve only websiteAddress and websitedescription.
	dnsRecords := convert.ByteToStringArray(dnsValue)
	if len(dnsRecords) <= indexOfWebsiteAddress {
		return "", "", fmt.Errorf("invalid website value: missing website address")
	}

	address := dnsRecords[indexOfWebsiteAddress]

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

// GetWebsiteMetadata retrieves candidate metadata of the website.
func GetWebsiteMetadata(client *node.Client, address string) ([]byte, error) {
	websiteMetadata, err := node.DatastoreEntry(client, address, convert.StringToBytes(metaKey))
	if err != nil {
		return nil, fmt.Errorf("reading key '%s' at '%s': %w", metaKey, address, err)
	}

	return websiteMetadata.CandidateValue, nil
}
