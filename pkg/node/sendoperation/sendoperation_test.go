package sendoperation

import (
	"encoding/base64"
	"testing"

	"github.com/massalabs/station/pkg/node/sendoperation/callsc"
)

func TestSerializeDeserializeCallSCMessage(t *testing.T) {
	// Create some example data
	expiry := uint64(123456)
	fee := uint64(789)
	address := "AU1MPDRXuR22mwYDFCeZUDgYjcTAF1co6xujx2X6ugoHeYeGY3B5"
	function := "exampleFunction"
	parameters := []byte("exampleParameters")
	gasLimit := uint64(1000000)
	coins := uint64(12345)

	// Create a new CallSC operation
	operation, err := callsc.New(address, function, parameters, gasLimit, coins)
	if err != nil {
		t.Fatalf("Failed to create CallSC operation: %v", err)
	}

	// Serialize the operation
	msg := message(expiry, fee, operation)
	msgB64 := base64.StdEncoding.EncodeToString(msg)

	// Simulate decoding and deserialization
	decodedMsg, err := DecodeMessage64(msgB64)
	if err != nil {
		t.Fatalf("Error decoding message: %v", err)
	}

	callSC, err := callsc.DecodeCallSCMessage(decodedMsg)
	if err != nil {
		t.Fatalf("Error decoding CallSC: %v", err)
	}

	// Verify the fields
	if callSC.Address != address {
		t.Errorf("Address mismatch: expected %s, got %s", address, callSC.Address)
	}

	if callSC.Function != function {
		t.Errorf("Function mismatch: expected %s, got %s", function, callSC.Function)
	}

	if string(callSC.Parameters) != string(parameters) {
		t.Errorf("Parameters mismatch: expected %v, got %v", parameters, callSC.Parameters)
	}

	if callSC.GasLimit != gasLimit {
		t.Errorf("GasLimit mismatch: expected %d, got %d", gasLimit, callSC.GasLimit)
	}

	if callSC.Coins != coins {
		t.Errorf("Coins mismatch: expected %d, got %d", coins, callSC.Coins)
	}
}
