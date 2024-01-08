package serializeaddress

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSerializeAddress(t *testing.T) {
	testCases := []struct {
		addr          string
		expectedBytes []byte
		expectError   bool
	}{
		{
			addr: "AU1MPDRXuR22mwYDFCeZUDgYjcTAF1co6xujx2X6ugoHeYeGY3B5",
			expectedBytes: []byte{
				0, 0, 46, 72, 55, 221, 83, 31, 169, 208, 41, 146, 210, 82, 27, 34, 114, 141, 159, 245, 209, 189, 40, 141, 126,
				123, 156, 223, 187, 205, 64, 236, 40, 184,
			},
			expectError: false,
		},
		{
			addr: "AS12YMz7NjyP3aeEWcSsiC58Hba8UxHapfGv7i4PmNMS2eKfmaqqC",
			expectedBytes: []byte{
				1, 0, 202, 232, 43, 43, 168, 202, 122, 146, 118, 233, 120, 40, 254, 229, 81, 255, 245, 51, 119, 228, 26, 142, 34,
				195, 43, 76, 3, 140, 20, 198, 15, 188,
			},
			expectError: false,
		},
		{
			addr:          "invalid-address",
			expectedBytes: nil,
			expectError:   true,
		},
	}

	for _, testCase := range testCases {
		bytes, err := SerializeAddress(testCase.addr)
		if testCase.expectError {
			if err == nil {
				t.Errorf("Expected an error but did not receive one for address %s", testCase.addr)
			}
		} else {
			if err != nil {
				t.Errorf("Received an unexpected error %v for address %s", err, testCase.addr)
			}

			if !bytesEqual(bytes, testCase.expectedBytes) {
				t.Errorf("Bytes mismatch for address %s: expected %v, but got %v", testCase.addr, testCase.expectedBytes, bytes)
			}
		}
	}
}

func bytesEqual(byteA, byteB []byte) bool {
	if len(byteA) != len(byteB) {
		return false
	}

	for i, v := range byteA {
		if v != byteB[i] {
			return false
		}
	}

	return true
}

func TestDeserializeAddress(t *testing.T) {
	testCases := []struct {
		versionedAddress []byte
		expectedAddress  string
	}{
		{
			versionedAddress: []byte{
				0, 0, 46, 72, 55, 221, 83, 31, 169, 208, 41, 146, 210, 82, 27, 34, 114, 141, 159, 245, 209, 189, 40, 141, 126,
				123, 156, 223, 187, 205, 64, 236, 40, 184,
			},
			expectedAddress: "AU1MPDRXuR22mwYDFCeZUDgYjcTAF1co6xujx2X6ugoHeYeGY3B5",
		},
	}

	for _, testCase := range testCases {
		address, err := DeserializeAddress(testCase.versionedAddress)

		require.NoError(t, err, "Received an unexpected error for versioned address %v", testCase.versionedAddress)
		assert.Equal(t, testCase.expectedAddress, address, "Address mismatch for versioned address %v",
			testCase.versionedAddress)
	}
}
