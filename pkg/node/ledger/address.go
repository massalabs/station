package ledger

import (
	"context"
	"fmt"
	"strings"

	"github.com/massalabs/thyra/pkg/node"
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

func Addresses(client *node.Client, addr []string) ([]Address, error) {
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

func KeysFiltered(client *node.Client, scAddress string, keyPrefix string) ([]string, error) {
	results, err := Addresses(client, []string{scAddress})
	if err != nil {
		return nil, fmt.Errorf("calling get_addresses with '%+v': %w", scAddress, err)
	}

	var filteredKeys []string

	for _, candidateDatastoreKey := range results[0].CandidateDatastoreKeys {

		if strings.Index(string(candidateDatastoreKey[4:]), keyPrefix) == 0 {
			filteredKeys = append(filteredKeys, string(candidateDatastoreKey[4:]))

		}
	}
	return filteredKeys, nil
}
