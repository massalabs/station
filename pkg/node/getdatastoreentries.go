package node

import (
	"context"
	"fmt"
	"strings"
)

type DatastoreEntryResponse struct {
	CandidateValue []byte `json:"candidate_value"`
	FinalValue     []byte `json:"final_value"`
}

type DatastoreEntry struct {
	Address string        `json:"address"`
	Key     JSONableSlice `json:"key"`
}

type JSONableSlice []byte

func (u JSONableSlice) MarshalJSON() ([]byte, error) {
	var result string
	if u == nil {
		result = "null"
	} else {
		result = strings.Join(strings.Fields(fmt.Sprintf("%d", u)), ",")
	}

	return []byte(result), nil
}

func NewDatastoreEntry(address string, key []byte) DatastoreEntry {
	return DatastoreEntry{
		Address: address,
		Key:     key,
	}
}

func FetchDatastoreEntry(client *Client, address string, key []byte) (*DatastoreEntryResponse, error) {
	entries := make([]DatastoreEntry, 1)
	entries[0] = NewDatastoreEntry(address, key)

	response, err := FetchDatastoreEntries(client, entries)
	if err != nil {
		return nil, err
	}

	return &response[0], nil
}

func ContractDatastoreEntries(client *Client, address string, keys [][]byte) ([]DatastoreEntryResponse, error) {
	entries := make([]DatastoreEntry, len(keys))

	for i, key := range keys {
		entries[i] = NewDatastoreEntry(address, key)
	}

	response, err := FetchDatastoreEntries(client, entries)
	if err != nil {
		return nil, fmt.Errorf("calling get_datastore_entries '%+v': %w", entries, err)
	}

	return response, nil
}

func FetchDatastoreEntries(client *Client, entries []DatastoreEntry) ([]DatastoreEntryResponse, error) {
	data := make([][]DatastoreEntry, 1)
	data[0] = entries

	response, err := client.RPCClient.Call(
		context.Background(),
		"get_datastore_entries",
		data,
	)
	if err != nil {
		return nil, fmt.Errorf("calling get_datastore_entries '%+v': %w", entries, err)
	}

	if response.Error != nil {
		return nil, response.Error
	}

	var entry []DatastoreEntryResponse

	err = response.GetObject(&entry)
	if err != nil {
		return nil, fmt.Errorf("parsing get_datastore_entries jsonrpc response '%+v': %w", response, err)
	}

	return entry, nil
}
