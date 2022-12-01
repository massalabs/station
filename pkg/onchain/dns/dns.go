package dns

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/massalabs/thyra/pkg/helper"
	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/node/base58"
	"github.com/massalabs/thyra/pkg/node/sendoperation"
	"github.com/massalabs/thyra/pkg/onchain"
	"github.com/massalabs/thyra/pkg/wallet"
)

const DNSRawAddress = "A1QxHhhi9crDJoEAaXRjkU2w3xsusBwpwGcAGpRRFAVUUuDWf7z"

func Resolve(client *node.Client, name string) (string, error) {
	const dnsPrefix = "record"

	entry, err := node.DatastoreEntry(client, DNSRawAddress, helper.StringToByteArray(dnsPrefix+name))
	if err != nil {
		return "", fmt.Errorf("calling node.DatastoreEntry with '%s' at '%s': %w", DNSRawAddress, dnsPrefix+name, err)
	}

	if len(entry.CandidateValue) == 0 {
		return "", errors.New("name not found")
	}

	return string(entry.CandidateValue), nil
}

type setApproval struct {
	Operator string `json:"operator"`
	Approved bool   `json:"approved"`
}

type setRecord struct {
	Name    []byte `json:"name"`
	Address string `json:"address"`
}

func SetRecord(client *node.Client, wallet wallet.Wallet, url string, smartContract string) (string, error) {
	addr, _, err := base58.VersionedCheckDecode(DNSRawAddress[1:])
	if err != nil {
		return "", fmt.Errorf("checking address '%s': %w", DNSRawAddress[1:], err)
	}

	rec := setRecord{
		Name:    helper.StringToByteArray(url),
		Address: smartContract,
	}

	param, err := json.Marshal(rec)
	if err != nil {
		return "", fmt.Errorf("marshalling '%+v': %w", rec, err)
	}

	result, err := onchain.CallFunction(client, wallet, addr, "setResolver", param, sendoperation.OneMassa)
	if err != nil {
		return "", fmt.Errorf("calling setResolver with '%+v' at '%s': %w", param, addr, err)
	}

	return result, nil
}

// DEAD CODE ??
func SetRecordManager(client *node.Client, wallet wallet.Wallet) (string, error) {
	addr, _, err := base58.VersionedCheckDecode(DNSRawAddress[1:])
	if err != nil {
		return "", fmt.Errorf("checking address '%s': %w", DNSRawAddress[1:], err)
	}

	// Set Resolver prepare data
	appr := &setApproval{
		Operator: wallet.Address,
		Approved: true,
	}

	param, err := json.Marshal(appr)
	if err != nil {
		return "", fmt.Errorf("marshalling '%+v': %w", appr, err)
	}

	result, err := onchain.CallFunction(client, wallet, addr, "setApprovalForAll", param, sendoperation.OneMassa)
	if err != nil {
		return "", fmt.Errorf("calling setApprovalForAll with '%+v' at '%s': %w", param, addr, err)
	}

	return result, nil
}
