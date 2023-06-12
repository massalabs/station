package websites

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/models"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/config"
	"github.com/massalabs/thyra/pkg/convert"
	"github.com/massalabs/thyra/pkg/node"
)

const (
	dateFormat          = "2006-01-02"
	metaKey             = "META"
	ownedPrefix         = "owned"
	ownerKey            = "owner"
	blackListKey        = "blackList"
	secondsToMilliCoeff = 1000
)

func NewRegistryHandler(config *config.AppConfig) operations.AllDomainsGetterHandler {
	return &registryHandler{config: config}
}

type registryHandler struct {
	config *config.AppConfig
}

func (h *registryHandler) Handle(params operations.AllDomainsGetterParams) middleware.Responder {
	results, err := Registry(*h.config)
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
func Registry(config config.AppConfig) ([]*models.Registry, error) {
	client := node.NewClient(config.NodeURL)

	websiteNames, err := filterEntriesToDisplay(config, client)
	if err != nil {
		return nil, fmt.Errorf("filtering keys to be displayed at '%s': %w", config.DNSAddress, err)
	}

	dnsValues, err := node.ContractDatastoreEntries(client, config.DNSAddress, websiteNames)
	if err != nil {
		return nil, fmt.Errorf("reading keys '%s' at '%s': %w", websiteNames, config.DNSAddress, err)
	}

	// in website name key, value are stored in this order -> website Address, website Owner Address,
	// website Description
	indexOfWebsiteAddress := 0
	indexOfWebsiteDescription := 2

	registry := make([]*models.Registry, len(dnsValues))

	for index := 0; index < len(dnsValues); index++ {
		websiteStorerAddress := convert.ByteToStringArray(dnsValues[index].CandidateValue)[indexOfWebsiteAddress]

		websiteDescription := convert.ByteToStringArray(dnsValues[index].CandidateValue)[indexOfWebsiteDescription]

		websiteMetadata, err := node.DatastoreEntry(client, websiteStorerAddress, convert.StringToBytes(metaKey))
		if err != nil {
			return nil, fmt.Errorf("reading key '%s' at '%s': %w", metaKey, websiteStorerAddress, err)
		}

		registry[index] = &models.Registry{
			Name:        convert.BytesToString(websiteNames[index]), // name of website : flappy.
			Address:     websiteStorerAddress,                       // website Address
			Description: websiteDescription,                         // website Description
			Metadata:    websiteMetadata.CandidateValue,             // website metadata.
		}
	}

	// sort website names with alphanumeric order.
	sort.Slice(registry, func(i, j int) bool {
		return registry[i].Name < registry[j].Name
	})

	return registry, nil
}

/*
The dns SC has 4 differents kinds of key :
-the website names
-keys owned concatenated with the owner's address
-a key blackList
-a owner key
we only want to keep the website names keys.
*/
func filterEntriesToDisplay(config config.AppConfig, client *node.Client) ([][]byte, error) {
	// we first remove the owned type keys
	keyList, err := node.FilterSCKeysByPrefix(client, config.DNSAddress, ownedPrefix, false)
	if err != nil {
		return nil, fmt.Errorf("fetching all keys without '%s' prefix at '%s': %w", ownedPrefix, config.DNSAddress, err)
	}

	// we then read the blacklisted websites
	blackListedWebsites, err := node.DatastoreEntry(client, config.DNSAddress, convert.StringToBytes(blackListKey))
	if err != nil {
		return nil, fmt.Errorf("reading entry '%s' prefix at '%s': %w", blackListKey, config.DNSAddress, err)
	}

	var keyListToRemove []string
	if !bytes.Equal(blackListedWebsites.CandidateValue, make([]byte, 0)) {
		keyListToRemove = strings.Split(convert.BytesToString(blackListedWebsites.CandidateValue), ",")
	}

	// we add the keys blackList and ownerKey to the list of key to be removed
	keyListToRemove = append(keyListToRemove, blackListKey, ownerKey)

	// we encode the list as a slice of byteArray
	keyListToRemoveAsArrayOfByteArray := convert.StringArrayToArrayOfByteArray(keyListToRemove)

	websiteNames := node.RemoveKeysFromKeyList(keyList, keyListToRemoveAsArrayOfByteArray)

	return websiteNames, nil
}
