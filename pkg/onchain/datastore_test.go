package onchain

import (
	"bytes"
	"testing"
)

func TestIsContractDeployDatastore(t *testing.T) {
	tests := []struct {
		name       string
		datastore  []datastoreEntry
		isDeploySC bool
	}{
		{
			name: "Valid contract deploy datastore",
			datastore: []datastoreEntry{
				{Key: []byte{0}, Value: []byte{1}},
				{Key: getContractByteCodeKey(), Value: []byte{1, 2, 3}},
				{Key: getArgsKey(), Value: []byte{4, 5, 6}},
				{Key: getCoinsKey(), Value: []byte{7, 8, 9}},
			},
			isDeploySC: true,
		},
		{
			name: "Invalid contract deploy datastore - wrong length",
			datastore: []datastoreEntry{
				{Key: []byte{0}, Value: []byte{1}},
				{Key: getContractByteCodeKey(), Value: []byte{1, 2, 3}},
			},
			isDeploySC: false,
		},
		{
			name: "Invalid contract deploy datastore - wrong smart contract number key",
			datastore: []datastoreEntry{
				{Key: []byte{1}, Value: []byte{1}},
				{Key: getContractByteCodeKey(), Value: []byte{1, 2, 3}},
				{Key: getArgsKey(), Value: []byte{4, 5, 6}},
				{Key: getCoinsKey(), Value: []byte{7, 8, 9}},
			},
			isDeploySC: false,
		},
		{
			name: "Invalid contract deploy datastore - wrong bytecode key",
			datastore: []datastoreEntry{
				{Key: []byte{0}, Value: []byte{1}},
				{Key: []byte("wrong key"), Value: []byte{1, 2, 3}},
				{Key: getArgsKey(), Value: []byte{4, 5, 6}},
				{Key: getCoinsKey(), Value: []byte{7, 8, 10}},
			},
			isDeploySC: false,
		},
		{
			name: "Invalid contract deploy datastore - wrong Args key",
			datastore: []datastoreEntry{
				{Key: []byte{0}, Value: []byte{1}},
				{Key: getContractByteCodeKey(), Value: []byte{1, 2, 3}},
				{Key: []byte("wrong key"), Value: []byte{4, 5, 6}},
				{Key: getCoinsKey(), Value: []byte{7, 8, 10}},
			},
			isDeploySC: false,
		},
		{
			name: "Invalid contract deploy datastore - wrong coins key",
			datastore: []datastoreEntry{
				{Key: []byte{0}, Value: []byte{1}},
				{Key: getContractByteCodeKey(), Value: []byte{1, 2, 3}},
				{Key: getArgsKey(), Value: []byte{4, 5, 6}},
				{Key: []byte("wrong key"), Value: []byte{7, 8, 10}},
			},
			isDeploySC: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isContractDeployDatastore(tt.datastore)
			if result != tt.isDeploySC {
				t.Errorf("isContractDeployDatastore() = %v, isDeploySC %v", result, tt.isDeploySC)
			}
		})
	}
}

func TestDeSerializeDatastore(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []datastoreEntry
	}{
		{
			name: "Valid datastore",
			input: func() []byte {
				datastore := []datastoreEntry{
					{Key: []byte{0}, Value: []byte{1}},
					{Key: []byte{1}, Value: []byte{2, 3}},
				}
				serialized, _ := SerializeDatastore(datastore)
				return serialized
			}(),
			expected: []datastoreEntry{
				{Key: []byte{0}, Value: []byte{1}},
				{Key: []byte{1}, Value: []byte{2, 3}},
			},
		},
		{
			name:     "Empty datastore",
			input:    []byte{},
			expected: nil,
		},
		{
			name: "Large datastore",
			input: func() []byte {
				datastore := []datastoreEntry{
					{Key: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, Value: []byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}},
				}
				serialized, _ := SerializeDatastore(datastore)
				return serialized
			}(),
			expected: []datastoreEntry{
				{Key: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, Value: []byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := DeSerializeDatastore(tt.input)
			if err != nil {
				t.Errorf("Got unexpected error = %v", err)
				return
			}

			if !equalDatastoreEntries(result, tt.expected) {
				t.Errorf("DeSerializeDatastore() = %+v, expected: %+v", result, tt.expected)
			}
		})
	}
}

func equalDatastoreEntries(a, b []datastoreEntry) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !bytes.Equal(a[i].Key, b[i].Key) || !bytes.Equal(a[i].Value, b[i].Value) {
			return false
		}
	}
	return true
}
