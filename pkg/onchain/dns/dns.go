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

const DNSRawAddress = "A1aNfHJ4CVHK4tW29jYcmx181zNWhf5GDyjqznV5HUrCsaSmCSD"

const bytesPerU32 = 4

func Resolve(client *node.Client, name string) (string, error) {
	const dnsPrefix = "record"

	entry, err := node.DatastoreEntry(client, DNSRawAddress, convert.EncodeStringUint32ToUTF8(dnsPrefix+name))
	if err != nil {
		return "", fmt.Errorf("calling node.DatastoreEntry with '%s' at '%s': %w", DNSRawAddress, dnsPrefix+name, err)
	}

	if len(entry.CandidateValue) == 0 {
		return "", errors.New("name not found")
	}
	// we remove from the address its header length expressed as a U32
	return string(entry.CandidateValue[bytesPerU32:]), nil
}

func SetRecord(client *node.Client, wallet wallet.Wallet, url string, smartContract string) (string, error) {
	addr, _, err := base58.VersionedCheckDecode(DNSRawAddress[1:])
	if err != nil {
		return "", fmt.Errorf("checking address '%s': %w", DNSRawAddress[1:], err)
	}

	// Set Resolver prepare data
	rec := []byte(convert.U32ToString(uint32(len(url))))
	rec = append(rec, []byte(url)...)
	rec = append(rec, convert.U32ToString(uint32(len(smartContract)))...)
	rec = append(rec, []byte(smartContract)...)

	result, err := onchain.CallFunction(client, wallet, addr, "setResolver", rec, sendoperation.OneMassa)
	if err != nil {
		return "", fmt.Errorf("calling setResolver with '%+v' at '%s': %w", rec, addr, err)
	}

	return result, nil
}
