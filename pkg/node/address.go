package node

import (
	"context"
	"fmt"
)

type Address struct {
	Address                string   `json:"address"`
	BlockDraws             []string `json:"block_draws"`
	BlocksCreated          []string `json:"blocks_created"`
	CandidateBalance       string   `json:"candidate_balance"`
	CandidateDatastoreKeys [][]byte `json:"candidate_datastore_keys"`
	FinalBalance           string   `json:"final_balance"`
	FinalDatastoreKeys     [][]byte `json:"final_datastore_keys"`
}

func Addresses(client *Client, addr []string) ([]Address, error) {
	response, err := client.RPCClient.Call(
		context.Background(),
		"get_addresses",
		[1][]string{addr})
	if err != nil {
		return nil, fmt.Errorf("calling get_addresses with '%+v': %w", [1][]string{addr}, err)
	}

	if response.Error != nil {
		return nil, response.Error
	}

	var content []Address
	err = response.GetObject(&content)
	if err != nil {
		return nil, fmt.Errorf("parsing get_addresses jsonrpc response '%+v': %w", response, err)
	}

	return content, nil
}
