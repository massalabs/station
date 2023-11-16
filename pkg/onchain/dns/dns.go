package dns

import (
	"errors"
	"fmt"
	"strings"

	"github.com/massalabs/station/int/config"
	"github.com/massalabs/station/pkg/convert"
	"github.com/massalabs/station/pkg/node"
	"github.com/massalabs/station/pkg/node/sendoperation"
	"github.com/massalabs/station/pkg/node/sendoperation/signer"
	"github.com/massalabs/station/pkg/onchain"
)

/*
This function fetch the address of the website storer associated with the name given in parameter
from the DNS smart contract and returns it.
*/
func Resolve(config config.NetworkInfos, client *node.Client, name string) (string, error) {
	entry, err := node.FetchDatastoreEntry(client, config.DNSAddress, convert.ToBytesWithPrefixLength(name))
	if err != nil {
		return "", fmt.Errorf("calling node.DatastoreEntry with '%s' at '%s': %w", config.DNSAddress, name, err)
	}

	if len(entry.CandidateValue) == 0 {
		return "", errors.New("name not found")
	}
	// entry.CandidateValue contains the website address + the owner address, we keep only the website address.
	return convert.ToStringArray(entry.CandidateValue)[0], nil
}

func SetRecord(
	config config.NetworkInfos,
	client *node.Client,
	nickname string,
	url string,
	description string,
	smartContract string,
	operationBatch sendoperation.OperationBatch,
) (string, error) {
	addr := config.DNSAddress

	// Set Resolver prepare data
	rec := convert.U32ToBytes(len(url))
	rec = append(rec, []byte(url)...)
	rec = append(rec, convert.U32ToBytes(len(smartContract))...)
	rec = append(rec, []byte(smartContract)...)
	rec = append(rec, convert.U32ToBytes(len(description))...)
	rec = append(rec, []byte(description)...)

	operationWithEventResponse, err := onchain.CallFunction(
		client,
		nickname,
		addr,
		"dns1_setResolver",
		rec,
		sendoperation.DefaultFee,
		sendoperation.DefaultGasLimitCallSC,
		sendoperation.OneMassa,
		sendoperation.DefaultExpiryInSlot,
		false,
		operationBatch,
		&signer.WalletPlugin{},
	)
	if err != nil {
		return "", fmt.Errorf("calling setResolver with '%+v' at '%s': %w", rec, addr, err)
	}

	event := operationWithEventResponse.Event

	if strings.HasPrefix(event, "ERROR") {
		return event, fmt.Errorf("calling setResolver failed with '%+v' at '%s': %s", rec, addr, event)
	}

	return operationWithEventResponse.Event, nil
}
