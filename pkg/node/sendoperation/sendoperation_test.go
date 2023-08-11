package sendoperation

import (
	"encoding/base64"
	"testing"

	"github.com/massalabs/station/pkg/node/sendoperation/callsc"
	"github.com/stretchr/testify/assert"
)

func TestSerializeDeserializeCallSCMessage(t *testing.T) {
	assert := assert.New(t)

	testestcaseases := []struct {
		expiry     uint64
		fee        uint64
		address    string
		function   string
		parameters []byte
		gasLimit   uint64
		coins      uint64
	}{
		{
			expiry:     uint64(123456),
			fee:        uint64(789),
			address:    "AU1MPDRXuR22mwYDFCeZUDgYjcTAF1co6xujx2X6ugoHeYeGY3B5",
			function:   "exampleFunction",
			parameters: []byte("exampleParameters"),
			gasLimit:   uint64(1000000),
			coins:      uint64(12345),
		},
	}

	for _, testcase := range testestcaseases {
		// Create a new CallSC operation
		operation, err := callsc.New(testcase.address, testcase.function, testcase.parameters,
			testcase.gasLimit, testcase.coins)
		assert.NoError(err, "Failed to create CallSC operation")

		// Serialize the operation
		msg := message(testcase.expiry, testcase.fee, operation)
		msgB64 := base64.StdEncoding.EncodeToString(msg)

		// Simulate decoding and deserialization
		decodedMsg, err := DecodeMessage64(msgB64)
		assert.NoError(err, "Error decoding message")

		callSC, err := callsc.DecodeMessage(decodedMsg)
		assert.NoError(err, "Error decoding CallSC")

		// Verify the fields
		assert.Equal(testcase.address, callSC.Address, "Address mismatestcaseh")
		assert.Equal(testcase.function, callSC.Function, "Function mismatestcaseh")
		assert.Equal(testcase.parameters, callSC.Parameters, "Parameters mismatestcaseh")
		assert.Equal(testcase.gasLimit, callSC.GasLimit, "GasLimit mismatestcaseh")
		assert.Equal(testcase.coins, callSC.Coins, "Coins mismatestcaseh")
	}
}
