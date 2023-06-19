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

// Helper function to fetch website values from DNS entry.
func GetWebsiteValues(dnsValue []byte) (string, string, error) {
	websiteValue := convert.ByteToStringArray(dnsValue)
	if len(websiteValue) <= indexOfWebsiteAddress {
		return "", "", fmt.Errorf("invalid website value: missing website address")
	}

	websiteStorerAddress := websiteValue[indexOfWebsiteAddress]

	websiteDescription := ""
	if len(websiteValue) > indexOfWebsiteDescription {
		websiteDescription = websiteValue[indexOfWebsiteDescription]
	}
	// Prevent XSS by escaping special characters in websiteDescription
	escapedDescription := template.HTMLEscapeString(websiteDescription)

	return websiteStorerAddress, escapedDescription, nil
}

// Helper function to retrieve website metadata.
func GetWebsiteMetadata(client *node.Client, websiteStorerAddress string) ([]byte, error) {
	websiteMetadata, err := node.DatastoreEntry(client, websiteStorerAddress, convert.StringToBytes(metaKey))
	if err != nil {
		return nil, fmt.Errorf("reading key '%s' at '%s': %w", metaKey, websiteStorerAddress, err)
	}

	return websiteMetadata.CandidateValue, nil
}
