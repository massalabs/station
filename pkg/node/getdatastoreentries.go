package node

import (
	"context"
	"fmt"
	"strings"
)

type getDatastoreEntries struct {
	Address string        `json:"address"`
	Key     JSONableSlice `json:"key"`
}

type DatastoreEntryResponse struct {
	CandidateValue []byte `json:"candidate_value"`
	FinalValue     []byte `json:"final_value"`
}

type DatastoreEntriesKeysAsString struct {
	Address string `json:"address"`
	Key     []byte `json:"key"`
}

type JSONableSlice []uint8

func (u JSONableSlice) MarshalJSON() ([]byte, error) {
	var result string
	if u == nil {
		result = "null"
	} else {
		result = strings.Join(strings.Fields(fmt.Sprintf("%d", u)), ",")
	}

	return []byte(result), nil
}

func DatastoreEntry(client *Client, address string, key []byte) (*DatastoreEntryResponse, error) {
	entries := []DatastoreEntriesKeysAsString{}

	entry := DatastoreEntriesKeysAsString{
		Address: address,
		Key:     key,
	}

	entries = append(entries, entry)

	response, err := DatastoreEntries(client, entries)
	if err != nil {
		return nil, err
	}

	return &response[0], nil
}

func ContractDatastoreEntries(client *Client, address string, keys []string) ([]DatastoreEntryResponse, error) {
	entries := []DatastoreEntriesKeysAsString{}

	for i := 0; i < len(keys); i++ {
		entry := DatastoreEntriesKeysAsString{
			Address: address,
			Key:     []byte(keys[i]),
		}

		entries = append(entries, entry)
	}

	response, err := DatastoreEntries(client, entries)
	if err != nil {
		return nil, fmt.Errorf("calling get_datastore_entries '%+v': %w", entries, err)
	}

	return response, nil
}

func DatastoreEntries(client *Client, params []DatastoreEntriesKeysAsString) ([]DatastoreEntryResponse, error) {
	entries := [][]getDatastoreEntries{

		{},
	}

	for i := 0; i < len(params); i++ {
		entry := getDatastoreEntries{
			Address: params[i].Address,
			Key:     params[i].Key,
		}

		entries[0] = append(entries[0], entry)
	}

	response, err := client.RPCClient.Call(
		context.Background(),
		"get_datastore_entries",
		entries,
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
