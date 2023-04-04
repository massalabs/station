package dns

import (
	"errors"
	"fmt"
	"os"

	"github.com/massalabs/thyra/pkg/convert"
	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/node/base58"
	"github.com/massalabs/thyra/pkg/node/sendoperation"
	"github.com/massalabs/thyra/pkg/onchain"
)

const EnvKey = "THYRA_DNS_ADDRESS"

func Address() string {
	return os.Getenv(EnvKey)
}

/*
This function fetch the address of the website storer associated with the name given in parameter
from the DNS smart contract and returns it.
*/
func Resolve(client *node.Client, name string) (string, error) {
	entry, err := node.DatastoreEntry(client, Address(), convert.StringToBytes(name))
	if err != nil {
		return "", fmt.Errorf("calling node.DatastoreEntry with '%s' at '%s': %w", Address(), name, err)
	}

	if len(entry.CandidateValue) == 0 {
		return "", errors.New("name not found")
	}
	// entry.CandidateValue contains the website address + the owner address, we keep only the website address.
	return convert.ByteToStringArray(entry.CandidateValue)[0], nil
}

func SetRecord(
	client *node.Client,
	nickname string,
	url string,
	smartContract string,
	operationBatch sendoperation.OperationBatch,
) (string, error) {
	addr, _, err := base58.VersionedCheckDecode(Address()[2:])
	if err != nil {
		return "", fmt.Errorf("checking address '%s': %w", Address()[2:], err)
	}

	// Set Resolver prepare data
	rec := convert.U32ToBytes(len(url))
	rec = append(rec, []byte(url)...)
	rec = append(rec, convert.U32ToBytes(len(smartContract))...)
	rec = append(rec, []byte(smartContract)...)

	operationWithEventResponse, err := onchain.CallFunction(
		client,
		nickname,
		addr,
		"setResolver",
		rec,
		sendoperation.OneMassa,
		operationBatch,
	)
	if err != nil {
		return "", fmt.Errorf("calling setResolver with '%+v' at '%s': %w", rec, addr, err)
	}

	return operationWithEventResponse.Event, nil
}

type MetaData struct {
	CreationTimeStamp   uint64
	LastUpdateTimestamp uint64
}

// FetchRecordMetaData returns the website meta data from the DNS samrt contract stored on the blockchain.
func FetchRecordMetaData(client *node.Client, websiteStorerAddress string) (MetaData, error) {
	data, err := node.DatastoreEntry(client, websiteStorerAddress, convert.StringToBytes("META"))
	if err != nil {
		return MetaData{}, fmt.Errorf("while getting meta data: %w", err)
	}

	creation := convert.BytesToU64(data.CandidateValue)
	update := uint64(0)

	if len(data.CandidateValue) == 2*convert.BytesPerUint64 {
		update = convert.BytesToU64(data.CandidateValue[convert.BytesPerUint64:])
	}

	return MetaData{
		CreationTimeStamp:   creation,
		LastUpdateTimestamp: update,
	}, nil
}
