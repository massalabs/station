package websites

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station/api/swagger/server/models"
	"github.com/massalabs/station/api/swagger/server/restapi/operations"
	"github.com/massalabs/station/int/config"
	"github.com/massalabs/station/pkg/convert"
	"github.com/massalabs/station/pkg/dnshelper"
	"github.com/massalabs/station/pkg/logger"
	"github.com/massalabs/station/pkg/node"
	"github.com/massalabs/station/pkg/onchain/website"
)

const (
	// Keys declared in the DNS smart contract.
	ownedPrefix  = "owned"
	ownerKey     = "OWNER"
	blackListKey = "blackList"
	// Other constants.
	faviconIcon = "favicon.ico"
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
// Registry fetches the registry data for the given AppConfig.
func Registry(config *config.AppConfig) ([]*models.Registry, error) {
	client := node.NewClient(config.NodeURL)

	websiteNames, err := filterEntriesToDisplay(*config, client)
	if err != nil {
		return nil, fmt.Errorf("failed to filter keys to be displayed at '%s': %w", config.DNSAddress, err)
	}

	dnsValues, err := node.ContractDatastoreEntries(client, config.DNSAddress, websiteNames)
	if err != nil {
		return nil, fmt.Errorf("failed to read keys '%s' at '%s': %w", websiteNames, config.DNSAddress, err)
	}

	registry := processDNSValues(client, dnsValues, websiteNames)

	sortRegistry(registry)

	return registry, nil
}

func processDNSValues(
	client *node.Client,
	dnsValues []node.DatastoreEntryResponse,
	websiteNames [][]byte,
) []*models.Registry {
	// Create buffered channels with the size of dnsValues slice
	registryChan := make(chan *models.Registry, len(dnsValues))
	errChan := make(chan error, len(dnsValues))
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(len(dnsValues))

	for index, dnsValue := range dnsValues {
		go func(index int, dnsValue node.DatastoreEntryResponse) {
			defer waitGroup.Done()

			processEntry(index, dnsValue, client, websiteNames, registryChan, errChan)
		}(index, dnsValue)
	}

	go func() {
		waitGroup.Wait()
		close(registryChan)
		close(errChan)
	}()

	return collectRegistryResults(registryChan, errChan)
}

func processEntry(index int,
	dnsValue node.DatastoreEntryResponse,
	client *node.Client, websiteNames [][]byte,
	registryChan chan<- *models.Registry, errChan chan<- error,
) {
	valueBytes := dnsValue.CandidateValue

	websiteStorerAddress, websiteDescription, err := dnshelper.AddressAndDescription(valueBytes)
	if err != nil {
		errChan <- err

		return
	}

	name := convert.ToString(websiteNames[index])

	registryChan <- &models.Registry{
		Name:        name,
		Address:     websiteStorerAddress,
		Description: websiteDescription,
		Favicon:     DNSRecordFavicon(name, websiteStorerAddress, client),
	}
}

func collectRegistryResults(registryChan <-chan *models.Registry, errChan <-chan error) []*models.Registry {
	registry := make([]*models.Registry, 0)

	for reg := range registryChan {
		registry = append(registry, reg)
	}

	if err := <-errChan; err != nil {
		// Handle the error if needed
		logger.Error(err)
	}

	return registry
}

func sortRegistry(registry []*models.Registry) {
	sort.Slice(registry, func(i, j int) bool {
		return registry[i].Name < registry[j].Name
	})
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
	blackListedWebsites, err := node.DatastoreEntry(
		client,
		config.DNSAddress,
		convert.ToBytesWithPrefixLength(blackListKey),
	)
	if err != nil {
		return nil, fmt.Errorf("reading entry '%s' prefix at '%s': %w", blackListKey, config.DNSAddress, err)
	}

	var keyListToRemove []string
	if !bytes.Equal(blackListedWebsites.CandidateValue, make([]byte, 0)) {
		keyListToRemove = strings.Split(convert.ToString(blackListedWebsites.CandidateValue), ",")
	}

	// we add the key blackList to the list of key to be removed
	keyListToRemove = append(keyListToRemove, blackListKey)

	// we encode the list as a slice of byteArray
	keyListToRemoveAsArrayOfByteArray := convert.StringArrayToArrayOfByteArray(keyListToRemove)

	// we add the key owner to the list of key to be removed
	keyListToRemoveAsArrayOfByteArray = append(
		keyListToRemoveAsArrayOfByteArray,
		convert.ToBytes(ownerKey),
	)

	websiteNames := node.RemoveKeysFromKeyList(keyList, keyListToRemoveAsArrayOfByteArray)

	return websiteNames, nil
}
