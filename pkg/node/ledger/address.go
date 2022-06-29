package ledger

import (
	"context"

	"github.com/massalabs/thyra/pkg/node"
)

type ledgerInfo struct {
	Datastore map[string][]byte
}

type Address struct {
	Info ledgerInfo `json:"candidate_sce_ledger_info"`
}

func Addresses(client *node.Client, addr []string) ([]Address, error) {
	response, err := client.RPCClient.Call(
		context.Background(),
		"get_addresses",
		[1][]string{addr})
	if err != nil {
		return nil, err
	}

	if response.Error != nil {
		return nil, response.Error
	}

	var content []Address

	err = response.GetObject(&content)
	if err != nil {
		return nil, err
	}

	return content, nil
}
