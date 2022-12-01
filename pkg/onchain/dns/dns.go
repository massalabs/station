package dns

import (
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

func SetRecord(client *node.Client, wallet wallet.Wallet, url string, smartContract string) (string, error) {
	addr, _, err := base58.VersionedCheckDecode(DNSRawAddress[1:])
	if err != nil {
		return "", fmt.Errorf("checking address '%s': %w", DNSRawAddress[1:], err)
	}

	// Set Resolver prepare data
	rec := []byte(helper.EncodeUint8ToUTF16String(uint32(len(url))))
	rec = append(rec, helper.StringToByteArray(url)...)
	rec = append(rec, helper.EncodeUint8ToUTF16String(uint32(len(smartContract)))...)
	rec = append(rec, helper.StringToByteArray(smartContract)...)

	result, err := onchain.CallFunction(client, wallet, addr, "setResolver", rec, sendoperation.OneMassa)
	if err != nil {
		return "", fmt.Errorf("calling setResolver with '%+v' at '%s': %w", rec, addr, err)
	}

	return result, nil
}
