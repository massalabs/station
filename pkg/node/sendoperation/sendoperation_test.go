package sendoperation

import (
	"encoding/base64"
	"testing"

	"github.com/massalabs/station/pkg/node/sendoperation/callsc"
	"github.com/massalabs/station/pkg/node/sendoperation/executesc"
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

func TestSerializeDeserializeExecuteSCMessage(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		expiry   uint64
		fee      uint64
		maxGas   uint64
		maxCoins uint64
	}{
		{
			expiry:   uint64(123456),
			fee:      uint64(789),
			maxGas:   uint64(100000),
			maxCoins: uint64(50000),
		},
	}

	for _, testCase := range testCases {
		// Create a new ExecuteSC operation
		operation := executesc.New(nil, testCase.maxGas, testCase.maxCoins, nil)

		// Serialize the operation
		msg := message(testCase.expiry, testCase.fee, operation)
		msgB64 := base64.StdEncoding.EncodeToString(msg)

		// Simulate decoding and deserialization
		decodedMsg, err := DecodeMessage64(msgB64)
		assert.NoError(err, "Error decoding message")

		executeSC, err := executesc.DecodeMessage(decodedMsg)
		assert.NoError(err, "Error decoding ExecuteSC")

		// Verify the fields
		assert.Equal(testCase.maxGas, executeSC.MaxGas, "MaxGas mismatch")
		assert.Equal(testCase.maxCoins, executeSC.MaxCoins, "MaxCoins mismatch")
	}
}
