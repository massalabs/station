package my

import (
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"

	"github.com/massalabs/station/api/swagger/server/models"
	"github.com/massalabs/station/int/config"
	"github.com/massalabs/station/pkg/convert"
	"github.com/massalabs/station/pkg/dnshelper"
	"github.com/massalabs/station/pkg/node"
	"github.com/massalabs/station/pkg/wallet"
)

//nolint:tagliatelle
type Domain struct {
	URL     string `json:"dnsName"`
	Address string `json:"address"`
}

/*
This function fetch the list of domain names owned by a user from the DNS smart contract
and returns it as an array of strings.
*/
func Domains(config config.NetworkInfos, client *node.Client, nickname string) ([]string, error) {
	const ownedPrefix = "owned"

	wallet, err := wallet.Fetch(nickname)
	if err != nil {
		return nil, fmt.Errorf("loading wallet '%s': %w", nickname, err)
	}

	names := []string{}

	ownerKey := convert.ToBytesWithPrefixLength(ownedPrefix + wallet.Address)

	rawNames, err := node.FetchDatastoreEntry(client, config.DNSAddress, ownerKey)
	if err != nil {
		return nil, fmt.Errorf("reading entry '%s' at '%s': %w", config.DNSAddress, ownedPrefix+wallet.Address, err)
	}

	if len(rawNames.CandidateValue) == 0 {
		return names, nil
	}

	names = strings.Split(convert.ToString(rawNames.CandidateValue), ",")

	if err != nil {
		return nil, fmt.Errorf("parsing json '%s': %w", rawNames.CandidateValue, err)
	}

	return names, nil
}

// GetWebsites retrieves information about websites given a DNSAddress, and domain names.
// It queries the DNS entries for the specified domain names, retrieves relevant data, checks chunk integrity,
// and returns a list of website information including contract address, name, description, and broken chunks.
func GetWebsites(config config.NetworkInfos, client *node.Client, domainNames []string) ([]*models.Websites, error) {
	keys := make([][]byte, len(domainNames))

	for i, domain := range domainNames {
		keys[i] = convert.ToBytesWithPrefixLength(domain)
	}

	// Retrieve DNS entries for the specified domain names
	dnsValues, err := node.ContractDatastoreEntries(client, config.DNSAddress, keys)
	if err != nil {
		return nil, fmt.Errorf("failed to read entries '%v': %w", keys, err)
	}

	// Store website information for each domain
	responses := make([]*models.Websites, len(domainNames))

	// Process each domain's DNS entry
	for index, domainName := range domainNames {
		// Extract contract address and website description from DNS entry
		websiteStorerAddress, websiteDescription, err := dnshelper.AddressAndDescription(dnsValues[index].CandidateValue)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve website values: %w", err)
		}

		// Check chunk integrity for the contract address
		missingChunks, err := getMissingChunkIds(client, websiteStorerAddress)
		if err != nil {
			return nil, fmt.Errorf("failed to check chunk integrity: %w", err)
		}

		// Create a response object with website information
		// Careful handle this in frontend:
		// 	BrokenChunks = [] => no missed chunks
		//  BrokenChunks = nil => website not stored on the blockchain

		response := &models.Websites{
			Address:      websiteStorerAddress,
			Name:         domainName,
			Description:  websiteDescription,
			BrokenChunks: missingChunks,
		}
		responses[index] = response
	}

	return responses, nil
}

// Check website chunks and store its ID in an array if one of them is broken.
func getMissingChunkIds(client *node.Client, address string) ([]string, error) {
	chunkNumberKey := "total_chunks"

	var missedChunks []string

	encodedNumberOfChunks, err := node.FetchDatastoreEntry(
		client,
		address,
		convert.ToBytesWithPrefixLength(chunkNumberKey),
	)
	if err != nil {
		return nil, fmt.Errorf("reading datastore entry '%s' at address '%s': %w", chunkNumberKey, address, err)
	}

	//nolint:gomnd
	if len(encodedNumberOfChunks.CandidateValue) < 8 /*sizeof uint64*/ {
		// If the key is not valid, it means that the website is not stored on the blockchain.
		return nil, nil
	}

	numberOfChunks := int(binary.LittleEndian.Uint64(encodedNumberOfChunks.CandidateValue))

	keys := make([][]byte, numberOfChunks)

	for i := 0; i < numberOfChunks; i++ {
		keys[i] = convert.ToBytesWithPrefixLength("massa_web_" + strconv.Itoa(i))
	}

	response, err := node.ContractDatastoreEntries(client, address, keys)
	if err != nil {
		return nil, fmt.Errorf("calling get_datastore_entries '%+v': %w", keys, err)
	}

	for i := 0; i < numberOfChunks; i++ {
		if string(response[i].CandidateValue) == "" {
			missedChunks = append(missedChunks, strconv.Itoa(i))
		}
	}

	return missedChunks, nil
}
