package websites

import (
	"fmt"
	"sort"
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

/*
This function fetch all websites data that are associated with the DNS
smart contract Thyra is connected to. Once this data has been fetched from the DNS and
the various website storer contracts, the function builds an array of Registry objects
and returns it to the frontend for display on the Registry page.
*/
func Registry(client *node.Client, candidateDatastoreKeys [][]byte) ([]*models.Registry, error) {
	// array of strings of website names : (recordflappy).
	recordKeysStrings, err := ledger.KeysFiltered(client, dns.DNSRawAddress, recordKey)
	if err != nil {
		return nil, fmt.Errorf("filtering keys with '%+v' failed : %w", recordKey, err)
	}
	// convert array of strings to array of [array of bytes]
	recordKeysBytes := make([][]byte, len(recordKeysStrings))
	for i, v := range recordKeysStrings {
		recordKeysBytes[i] = convert.StringToBytes(v)
	}
	// retrieve the records owners values : addresses who own the websites.
	recordResult, err := node.ContractDatastoreEntries(client, dns.DNSRawAddress, recordKeysBytes)
	if err != nil {
		return nil, fmt.Errorf("searching Owners of records (addresses) failed : %w", err)
	}

	var websiteStorers []node.DatastoreEntriesKeysAsString

	for _, record := range recordResult {
		if wallet.CheckAddress(convert.BytesToString(record.CandidateValue)) {
			websiteStorerKey := node.DatastoreEntriesKeysAsString{
				Address: convert.BytesToString(record.CandidateValue),
				Key:     convert.StringToBytes(metaKey + string(record.CandidateValue)),
			}

			websiteStorers = append(websiteStorers, websiteStorerKey)
		}
	}

	websitesMetadata, err := node.DatastoreEntries(client, websiteStorers)
	if err != nil {
		return nil, fmt.Errorf("metadata reaching on dnsContractStorers failed : %w", err)
	}

	registry := make([]*models.Registry, len(websitesMetadata))

	for index := 0; index < len(websitesMetadata); index++ {
		registry[index] = &models.Registry{
			Name:     strings.Split(recordKeysStrings[index], recordKey)[1], // name of website : flappy.
			Address:  websiteStorers[index].Address,                         // owner of Website Address.
			Metadata: websitesMetadata[index].CandidateValue,                // website metadata.
		}
	}
	// sort website names with alphanumeric order.
	sort.Slice(registry, func(i, j int) bool {
		return registry[i].Name < registry[j].Name
	})

	return registry, nil
}
