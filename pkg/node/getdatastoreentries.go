package node

import (
	"context"
	"errors"
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

func DatastoreEntry(client *Client, address string, key string) (*DatastoreEntryResponse, error) {
	response, err := client.RPCClient.Call(
		context.Background(),
		"get_datastore_entries",
		[][]getDatastoreEntries{
			{
				getDatastoreEntries{
					Address: address,
					Key:     []byte(key),
				},
			},
		})

	if err != nil {
		return nil, err
	}

	if response.Error != nil {
		return nil, response.Error
	}

	var entry []DatastoreEntryResponse

	err = response.GetObject(&entry)
	if err != nil {
		return nil, err
	}

	if len(entry) < 1 {
		return nil, errors.New("no entry")
	}

	return &entry[0], nil
}

type GetDatastoreEntriesString struct {
	Address string `json:"address"`
	Key     string `json:"key"`
}

func DatastoreEntries(client *Client, params []GetDatastoreEntriesString) (*[]DatastoreEntryResponse, error) {
	entries := [][]getDatastoreEntries{
		{},
	}
	for i := 0; i < len(params); i++ {
		entry := getDatastoreEntries{
			Address: params[i].Address,
			Key:     []byte(params[i].Key),
		}
		entries[0] = append(entries[0], entry)
	}

	response, err := client.RPCClient.Call(
		context.Background(),
		"get_datastore_entries",
		entries,
	)

	if err != nil {
		return nil, err
	}

	if response.Error != nil {
		return nil, response.Error
	}

	var entry []DatastoreEntryResponse

	err = response.GetObject(&entry)
	if err != nil {
		return nil, err
	}

	if len(entry) < 1 {
		return nil, errors.New("no entry")
	}

	return &entry, nil
}
