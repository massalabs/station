package ledger

import (
	"context"
	"fmt"

	"github.com/massalabs/thyra/pkg/node"
)

type ledgerInfo struct {
	Datastore map[string][]byte
}

type Address struct {
	CandidateSCELedgerInfo ledgerInfo `json:"candidate_sce_ledger_info"`
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
