package websites

import (
	"fmt"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/models"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/convert"
	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/node/ledger"
	"github.com/massalabs/thyra/pkg/onchain/dns"
	"github.com/massalabs/thyra/pkg/wallet"
)

const (
	dateFormat          = "2006-01-02"
	recordKey           = "record"
	metaKey             = "META"
	secondsToMilliCoeff = 1000
)

//nolint:nolintlint,ireturn
func RegistryHandler(params operations.AllDomainsGetterParams) middleware.Responder {
	client := node.NewDefaultClient()

	addressesResult, err := ledger.Addresses(client, []string{dns.DNSRawAddress})
	if err != nil {
		return operations.NewMyDomainsGetterInternalServerError().
			WithPayload(
				&models.Error{
					Code:    errorCodeGetRegistry,
					Message: err.Error(),
				})
	}

	results, err := Registry(client, addressesResult[0].CandidateDatastoreKeys)
	if err != nil {
		return operations.NewMyDomainsGetterInternalServerError().
			WithPayload(
				&models.Error{
					Code:    errorCodeGetRegistry,
					Message: err.Error(),
				})
	}

	return operations.NewAllDomainsGetterOK().WithPayload(results)
}

func Registry(client *node.Client, candidateDatastoreKeys [][]byte) ([]*models.Registry, error) {
	recordKeysStrings, err := ledger.KeysFiltered(client, dns.DNSRawAddress, recordKey) // array of strings of website names : (recordflappy)
	if err != nil {
		return nil, fmt.Errorf("filtering keys with '%+v' failed : %w", recordKey, err)
	}
	// convert array of strings to array of [array of bytes]
	recordKeyBytes := make([][]byte, len(recordKeysStrings))
	for i, v := range recordKeysStrings {
		recordKeyBytes[i] = convert.EncodeStringUint32ToUTF8(v)
	}
	// retrieve the records owners values : addresses who own the websites
	recordResult, err := node.ContractDatastoreEntries(client, dns.DNSRawAddress, recordKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("searching Owners of records (addresses) failed : %w", err)
	}

	var metadataKeys []node.DatastoreEntriesKeysAsString

	for _, record := range recordResult {
		if wallet.CheckAddress(convert.DecodeStringUTF8ToUint32(record.CandidateValue)) {
			metadataKey := node.DatastoreEntriesKeysAsString{
				Address: convert.DecodeStringUTF8ToUint32(record.CandidateValue),
				Key:     convert.EncodeStringUint32ToUTF8(metaKey + string(record.CandidateValue)),
			}

			metadataKeys = append(metadataKeys, metadataKey)
		}
	}

	metadatas, err := node.DatastoreEntries(client, metadataKeys)
	if err != nil {
		return nil, fmt.Errorf("metadata reaching on dnsContractStorers failed : %w", err)
	}

	registryResult := make([]*models.Registry, len(metadatas))

	// Here we have to switch to string to display to the front
	for index := 0; index < len(metadatas); index++ {
		registryResult[index] = &models.Registry{
			Name:     strings.Split(recordKeysStrings[index], recordKey)[1], // name of website : flappy
			Address:  metadataKeys[index].Address[4:],                       // Owner of Website Address
			Metadata: metadatas[index].CandidateValue,                       // TimeStamp : Date of Upload
		}
	}

	return registryResult, nil
}
