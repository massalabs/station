package my

import (
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"

	"github.com/massalabs/thyra/api/swagger/server/models"
	"github.com/massalabs/thyra/pkg/config"
	"github.com/massalabs/thyra/pkg/convert"
	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/wallet"
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
func Domains(config config.AppConfig, client *node.Client, nickname string) ([]string, error) {
	const ownedPrefix = "owned"

	wallet, err := wallet.Fetch(nickname)
	if err != nil {
		return nil, fmt.Errorf("loading wallet '%s': %w", nickname, err)
	}

	names := []string{}

	ownerKey := convert.StringToBytes(ownedPrefix + wallet.Address)

	rawNames, err := node.DatastoreEntry(client, config.DNSAddress, ownerKey)
	if err != nil {
		return nil, fmt.Errorf("reading entry '%s' at '%s': %w", config.DNSAddress, ownedPrefix+wallet.Address, err)
	}

	if len(rawNames.CandidateValue) == 0 {
		return names, nil
	}

	names = strings.Split(convert.BytesToString(rawNames.CandidateValue), ",")

	if err != nil {
		return nil, fmt.Errorf("parsing json '%s': %w", rawNames.CandidateValue, err)
	}

	return names, nil
}

func Websites(config config.AppConfig, client *node.Client, domainNames []string) ([]*models.Websites, error) {
	params := []node.DatastoreEntriesKeys{}

	for i := 0; i < len(domainNames); i++ {
		param := node.DatastoreEntriesKeys{
			Address: config.DNSAddress,
			Key:     convert.StringToBytes(domainNames[i]),
		}
		params = append(params, param)
	}

	responses := []*models.Websites{}

	contractAddresses, err := node.DatastoreEntries(client, params)
	if err != nil {
		return nil, fmt.Errorf("reading entries'%s': %w", params, err)
	}

	for i := 0; i < len(domainNames); i++ { //nolint:varnamelen
		contractAddressesIndex := 0
		websiteDescriptionIndex := 2

		contractAddress := convert.ByteToStringArray(contractAddresses[i].CandidateValue)[contractAddressesIndex]

		websiteDescription := convert.ByteToStringArray(contractAddresses[i].CandidateValue)[websiteDescriptionIndex]

		missingChunks, err := getMissingChunkIds(client, contractAddress)
		if err != nil {
			return nil, fmt.Errorf("checking chunk integrity: %w", err)
		}

		response := models.Websites{
			Address:      contractAddress,
			Name:         domainNames[i],
			Description:  websiteDescription,
			BrokenChunks: missingChunks,
		}
		responses = append(responses, &response)
	}

	return responses, nil
}

// Check website chunks and store its ID in an array if one of them is broken.
func getMissingChunkIds(client *node.Client, address string) ([]string, error) {
	chunkNumberKey := "total_chunks"

	var missedChunks []string

	encodedNumberOfChunks, err := node.DatastoreEntry(client, address, convert.StringToBytes(chunkNumberKey))
	if err != nil {
		return nil, fmt.Errorf("reading datastore entry '%s' at address '%s': %w", chunkNumberKey, address, err)
	}

	numberOfChunks := int(binary.LittleEndian.Uint64(encodedNumberOfChunks.CandidateValue))

	entries := []node.DatastoreEntriesKeys{}

	for i := 0; i < numberOfChunks; i++ {
		entry := node.DatastoreEntriesKeys{
			Address: address,
			Key:     convert.StringToBytes("massa_web_" + strconv.Itoa(i)),
		}
		entries = append(entries, entry)
	}

	response, err := node.DatastoreEntries(client, entries)
	if err != nil {
		return nil, fmt.Errorf("calling get_datastore_entries '%+v': %w", entries, err)
	}

	for i := 0; i < numberOfChunks; i++ {
		if string(response[i].CandidateValue) == "" {
			missedChunks = append(missedChunks, strconv.Itoa(i))
		}
	}

	return missedChunks, nil
}
