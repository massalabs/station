package websites

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station/api/swagger/server/models"
	"github.com/massalabs/station/api/swagger/server/restapi/operations"
	"github.com/massalabs/station/pkg/config"
	"github.com/massalabs/station/pkg/convert"
	"github.com/massalabs/station/pkg/dnshelper"
	"github.com/massalabs/station/pkg/node"
	"github.com/massalabs/station/pkg/onchain/website"
)

const (
	dateFormat          = "2006-01-02"
	ownedPrefix         = "owned"
	ownerKey            = "owner"
	blackListKey        = "blackList"
	secondsToMilliCoeff = 1000
	faviconIcon         = "favicon.ico"
)

func NewRegistryHandler(config *config.AppConfig) operations.AllDomainsGetterHandler {
	return &registryHandler{config: config}
}

type registryHandler struct {
	config *config.AppConfig
}

func (h *registryHandler) Handle(_ operations.AllDomainsGetterParams) middleware.Responder {
	results, err := Registry(h.config)
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
Registry fetches all websites data that are associated with the DNS
smart contract Massa Station is connected to. Once this data has been fetched from the DNS and
the various website storer contracts, the function builds an array of Registry objects
and returns it to the frontend for display on the Registry page.
*/
func Registry(config *config.AppConfig) ([]*models.Registry, error) {
	client := node.NewClient(config.NodeURL)

	// Fetch website names to display
	websiteNames, err := filterEntriesToDisplay(*config, client)
	if err != nil {
		return nil, fmt.Errorf("filtering keys to be displayed at '%s': %w", config.DNSAddress, err)
	}

	// Fetch DNS values of the website names: contractAddress, Description.
	dnsValues, err := node.ContractDatastoreEntries(client, config.DNSAddress, websiteNames)
	if err != nil {
		return nil, fmt.Errorf("reading keys '%s' at '%s': %w", websiteNames, config.DNSAddress, err)
	}

	// Pre-allocate registry with the estimated capacity
	registry := make([]*models.Registry, 0, len(dnsValues))

	// Iterate over the DNS values to populate Registry objects with website values
	for index, dnsValue := range dnsValues {
		websiteStorerAddress, websiteDescription, err := dnshelper.AddressAndDescription(dnsValue.CandidateValue)
		if err != nil {
			// Skip invalid entry and continue to the next iteration
			continue
		}

		websiteMetadata, err := dnshelper.GetWebsiteMetadata(client, websiteStorerAddress)
		if err != nil {
			// Skip invalid entry and continue to the next iteration
			continue
		}

		name := convert.BytesToString(websiteNames[index])

		registry = append(registry, &models.Registry{
			Name:        name,
			Address:     websiteStorerAddress,
			Description: websiteDescription,
			Metadata:    websiteMetadata,
			Favicon:     DNSRecordFavicon(name, websiteStorerAddress, client),
		})
	}

	// Sort website names with alphanumeric order
	sort.Slice(registry, func(i, j int) bool {
		return registry[i].Name < registry[j].Name
	})

	return registry, nil
}

func DNSRecordFavicon(name, websiteStorerAddress string, client *node.Client) string {
	body, err := website.Fetch(client, websiteStorerAddress, faviconIcon)
	if err != nil || len(body) == 0 {
		return ""
	}

	return "https://" + name + ".massa/" + faviconIcon
}

/*
The dns SC has 4 different kinds of key:
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
