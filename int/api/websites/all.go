package websites

import (
	"bytes"
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
)

const (
	dateFormat          = "2006-01-02"
	metaKey             = "META"
	ownedPrefix         = "owned"
	blackListKey        = "blackList"
	secondsToMilliCoeff = 1000
)

//nolint:nolintlint,ireturn
func RegistryHandler(params operations.AllDomainsGetterParams) middleware.Responder {
	client := node.NewDefaultClient()

	results, err := Registry(client)
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
func Registry(client *node.Client) ([]*models.Registry, error) {
	websiteNames, err := filterEntriesToDisplay(client)
	if err != nil {
		return nil, fmt.Errorf("filtering keys to be displayed at '%s': %w", dns.DNSRawAddress, err)
	}

	dnsValues, err := node.ContractDatastoreEntries(client, dns.Address(), websiteNames)
	if err != nil {
		return nil, fmt.Errorf("reading keys '%s' at '%s': %w", websiteNames, dns.Address(), err)
	}

	// in website name key, value are stored in this order -> website Address, website Owner Address
	indexOfWebsiteAddress := 0

	registry := make([]*models.Registry, len(dnsValues))

	for index := 0; index < len(dnsValues); index++ {
		websiteStorerAddress := convert.ByteToStringArray(dnsValues[index].CandidateValue)[indexOfWebsiteAddress]

		websiteMetadata, err := node.DatastoreEntry(client, websiteStorerAddress, convert.StringToBytes(metaKey))
		if err != nil {
			return nil, fmt.Errorf("reading key '%s' at '%s': %w", metaKey, websiteStorerAddress, err)
		}

		registry[index] = &models.Registry{
			Name:     convert.BytesToString(websiteNames[index]), // name of website : flappy.
			Address:  websiteStorerAddress,                       // website Address
			Metadata: websiteMetadata.CandidateValue,             // website metadata.
		}
	}

	// sort website names with alphanumeric order.
	sort.Slice(registry, func(i, j int) bool {
		return registry[i].Name < registry[j].Name
	})

	return registry, nil
}

/*
The dns SC has 3 different kinds of key :
-the website names
-keys owned concatenated with the owner's address
-a key blackList
we only want to keep the website names keys.
*/
func filterEntriesToDisplay(client *node.Client) ([][]byte, error) {
	// we first remove the owned type keys
	keyList, err := ledger.FilterSCKeysByPrefix(client, dns.DNSRawAddress, ownedPrefix, false)
	if err != nil {
		return nil, fmt.Errorf("fetching all keys without '%s' prefix at '%s': %w", ownedPrefix, dns.DNSRawAddress, err)
	}

	// we then read the blacklisted websites
	blackListedWebsites, err := node.DatastoreEntry(client, dns.DNSRawAddress, convert.StringToBytes(blackListKey))
	if err != nil {
		return nil, fmt.Errorf("reading entry '%s' prefix at '%s': %w", blackListKey, dns.DNSRawAddress, err)
	}

	var keyListToRemove []string
	if !bytes.Equal(blackListedWebsites.CandidateValue, make([]byte, 0)) {
		keyListToRemove = strings.Split(convert.BytesToString(blackListedWebsites.CandidateValue), ",")
	}

	// we add the key blackList to the list of key to be removed
	keyListToRemove = append(keyListToRemove, blackListKey)

	// we encode the list as a slice of byteArray
	keyListToRemoveAsArrayOfByteArray := convert.StringArrayToArrayOfByteArray(keyListToRemove)

	websiteNames := ledger.RemoveKeysFromKeyList(keyList, keyListToRemoveAsArrayOfByteArray)

	return websiteNames, nil
}
