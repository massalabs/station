package websites

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/models"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/node/ledger"
	"github.com/massalabs/thyra/pkg/onchain/dns"
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

type dateOnChain struct {
	CreateDate int64 `json:"create_date"`
	UpdateDate int64 `json:"update_date"`
}

func Registry(client *node.Client, candidateDatastoreKeys [][]byte) ([]*models.Registry, error) {
	recordKeys, err := ledger.KeysFiltered(client, dns.DNSRawAddress, recordKey)
	if err != nil {
		return nil, fmt.Errorf("filtering keys with '%+v' failed : %w", recordKey, err)
	}

	recordResult, err := node.DatastoreEntriesOnSameContract(client, dns.DNSRawAddress, recordKeys)
	if err != nil {
		return nil, fmt.Errorf("searching recordAddress failed : %w", err)
	}

	var metadataKeys []node.DatastoreEntriesKeysAsString
	for index, record := range recordResult {
		metadataKeys[index] = node.DatastoreEntriesKeysAsString{
			Address: string(record.CandidateValue), Key: metaKey,
		}
	}

	metadatas, err := node.DatastoreEntries(client, metadataKeys)
	if err != nil {
		return nil, fmt.Errorf("metadata reaching on dnsContractStorers failed : %w", err)
	}

	var dates []dateOnChain

	for index, metadata := range metadatas {
		var date dateOnChain

		_ = json.Unmarshal(metadata.CandidateValue, &date)
		dates[index] = date
	}

	var registryResult []*models.Registry

	for index, date := range dates {
		registryResult[index] = &models.Registry{
			Name:      strings.Split(recordKeys[index], recordKey)[1],
			Address:   metadataKeys[index].Address,
			CreatedAt: time.Unix(date.CreateDate/secondsToMilliCoeff, 0).Format(dateFormat),
			UpdatedAt: time.Unix(date.UpdateDate/secondsToMilliCoeff, 0).Format(dateFormat),
		}
	}

	return registryResult, nil
}
