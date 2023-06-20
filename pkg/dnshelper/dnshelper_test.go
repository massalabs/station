package dnshelper

import (
	"testing"

	"github.com/massalabs/station/pkg/convert"
)

func TestAddressAndDescription(t *testing.T) {
	websiteAddressStr := ""
	ownerAddressStr := ""
	descriptionStr := "<script>alert('XSS');</script>"

	websiteAddressBytes := convert.StringToBytes(websiteAddressStr)
	ownerAddressBytes := convert.StringToBytes(ownerAddressStr)
	descriptionBytes := convert.StringToBytes(descriptionStr)

	dnsValue := append(append(websiteAddressBytes, ownerAddressBytes...), (descriptionBytes)...)
	expectedDescription := "&lt;script&gt;alert(&#39;XSS&#39;);&lt;/script&gt;"

	_, description, err := AddressAndDescription(dnsValue)
	if err != nil {
		t.Errorf("Unexpected error while calling AddressAndDescription: %s", err)

		return
	}

	if description != expectedDescription {
		t.Errorf("Description mismatch. Expected: %s, Got: %s", expectedDescription, description)
	}
}
