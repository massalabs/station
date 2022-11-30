package websites

import (
	"fmt"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/models"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/helper"
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
	recordKeys, err := ledger.KeysFiltered(client, dns.DNSRawAddress, recordKey)
	// helper.StringToByteArray(recordKeys)
	if err != nil {
		return nil, fmt.Errorf("filtering keys with '%+v' failed : %w", recordKey, err)
	}

	recordResult, err := node.ContractDatastoreEntries(client, dns.DNSRawAddress, recordKeys)
	if err != nil {
		return nil, fmt.Errorf("searching recordAddress failed : %w", err)
	}

	var metadataKeys []node.DatastoreEntriesKeysAsString

	for _, record := range recordResult {
		if wallet.AddressChecker(string(record.CandidateValue)) {
			metadataKeys = append(metadataKeys, node.DatastoreEntriesKeysAsString{
				Address: string(record.CandidateValue), Key: helper.StringToByteArray(metaKey),
			})
		}
	}

	metadatas, err := node.DatastoreEntries(client, metadataKeys)
	if err != nil {
		return nil, fmt.Errorf("metadata reaching on dnsContractStorers failed : %w", err)
	}

	registryResult := make([]*models.Registry, len(metadatas))

	for index := 0; index < len(metadatas); index++ {
		registryResult[index] = &models.Registry{
			Name:     strings.Split(recordKeys[index], recordKey)[1],
			Address:  metadataKeys[index].Address,
			Metadata: metadatas[index].CandidateValue,
		}
	}

	return registryResult, nil
}
