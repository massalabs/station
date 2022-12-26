package dns

import (
	"errors"
	"fmt"

	"github.com/massalabs/thyra/pkg/convert"
	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/node/base58"
	"github.com/massalabs/thyra/pkg/node/sendoperation"
	"github.com/massalabs/thyra/pkg/onchain"
	"github.com/massalabs/thyra/pkg/wallet"
)

const DNSRawAddress = "A1HEaAoZhuBGKFGJmpEwwrcy25EuzcNobxtYiFAq6R5hfNsH7j4"

/*
This function fetch the address of the website storer associated with the name given in parameter
from the DNS smart contract and returns it.
*/
func Resolve(client *node.Client, name string) (string, error) {
	entry, err := node.DatastoreEntry(client, DNSRawAddress, convert.StringToBytes(name))
	if err != nil {
		return "", fmt.Errorf("calling node.DatastoreEntry with '%s' at '%s': %w", DNSRawAddress, name, err)
	}

	if len(entry.CandidateValue) == 0 {
		return "", errors.New("name not found")
	}
	// in the name key is stored the storer Address + the owner Address, we keep only the storer Address
	return convert.ByteToStringArray(entry.CandidateValue)[0], nil
}

func SetRecord(client *node.Client, wallet wallet.Wallet, url string, smartContract string) (string, error) {
	addr, _, err := base58.VersionedCheckDecode(DNSRawAddress[1:])
	if err != nil {
		return "", fmt.Errorf("checking address '%s': %w", DNSRawAddress[1:], err)
	}

	// Set Resolver prepare data
	rec := convert.U32ToBytes(len(url))
	rec = append(rec, []byte(url)...)
	rec = append(rec, convert.U32ToBytes(len(smartContract))...)
	rec = append(rec, []byte(smartContract)...)

	result, err := onchain.CallFunction(client, wallet, addr, "setResolver", rec, sendoperation.OneMassa)
	if err != nil {
		return "", fmt.Errorf("calling setResolver with '%+v' at '%s': %w", rec, addr, err)
	}

	return result, nil
}
