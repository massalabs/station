package dns

import (
	"encoding/json"

	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/node/base58"
	"github.com/massalabs/thyra/pkg/onchain"
	"github.com/massalabs/thyra/pkg/wallet"
)

type setApproval struct {
	Operator string `json:"operator"`
	Approved bool   `json:"approved"`
}

type setRecord struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

const dnsRawAddress = "A12ew8eiCS7wnY8SkUdwBgDkdD5qwmbJgkJvYLCvVjWWdoFJJLvW"

func SetRecord(client *node.Client, wallet wallet.Wallet, url string, smartContract string) (string, error) {
	addr, _, err := base58.VersionedCheckDecode(dnsRawAddress[1:])
	if err != nil {
		return "", err
	}

	rec := setRecord{
		Name:    url,
		Address: smartContract,
	}

	param, err := json.Marshal(rec)
	if err != nil {
		return "", err
	}

	return onchain.CallFunction(client, wallet, addr, "setResolver", param)
}

func SetRecordManager(client *node.Client, wallet wallet.Wallet) (string, error) {
	addr, _, err := base58.VersionedCheckDecode(dnsRawAddress[1:])
	if err != nil {
		return "", err
	}

	// Set Resolver prepare data
	appr := &setApproval{
		Operator: wallet.Address,
		Approved: true,
	}

	param, err := json.Marshal(appr)
	if err != nil {
		return "", err
	}

	return onchain.CallFunction(client, wallet, addr, "setApprovalForAll", param)
}
