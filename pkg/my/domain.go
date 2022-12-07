package my

import (
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"

	"github.com/massalabs/thyra/api/swagger/server/models"
	"github.com/massalabs/thyra/pkg/convert"
	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/onchain/dns"
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
func Domains(client *node.Client, nickname string) ([]string, error) {
	const ownedPrefix = "owned"

	wallet, err := wallet.Load(nickname)
	if err != nil {
		return nil, fmt.Errorf("loading wallet '%s': %w", nickname, err)
	}

	domainsList := []string{}

	userDomainListKey := convert.EncodeStringToByteArray(ownedPrefix + wallet.Address)

	encodedUserDomainsList, err := node.DatastoreEntry(client, dns.DNSRawAddress, userDomainListKey)
	if err != nil {
		return nil, fmt.Errorf("reading entry '%s' at '%s': %w", dns.DNSRawAddress, ownedPrefix+wallet.Address, err)
	}

	if len(encodedUserDomainsList.CandidateValue) == 0 {
		return domainsList, nil
	}

	domainsList = strings.Split(convert.RemoveStringEncodingPrefix(encodedUserDomainsList.CandidateValue), ",")

	if err != nil {
		return nil, fmt.Errorf("parsing json '%s': %w", encodedUserDomainsList.CandidateValue, err)
	}

	return domainsList, nil
}

func Websites(client *node.Client, domainNames []string) ([]*models.Websites, error) {
	const recordPrefix = "record"

	params := []node.DatastoreEntriesKeysAsString{}

	for i := 0; i < len(domainNames); i++ {
		param := node.DatastoreEntriesKeysAsString{
			Address: dns.DNSRawAddress,
			Key:     convert.EncodeStringToByteArray(recordPrefix + domainNames[i]),
		}
		params = append(params, param)
	}

	responses := []*models.Websites{}

	contractAddresses, err := node.DatastoreEntries(client, params)
	if err != nil {
		return nil, fmt.Errorf("reading entries'%s': %w", params, err)
	}

	for i := 0; i < len(domainNames); i++ { //nolint:varnamelen
		contractAddress := convert.RemoveStringEncodingPrefix(contractAddresses[i].CandidateValue)

		missingChunks, err := getMissingChunkIds(client, contractAddress)
		if err != nil {
			return nil, fmt.Errorf("checking chunk integrity: %w", err)
		}

		response := models.Websites{
			Address:      contractAddress,
			Name:         domainNames[i],
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

	encodedNumberOfChunks, err := node.DatastoreEntry(client, address, convert.EncodeStringToByteArray(chunkNumberKey))
	if err != nil {
		return nil, fmt.Errorf("reading datastore entry '%s' at address '%s': %w", chunkNumberKey, address, err)
	}

	numberOfChunks := int(binary.LittleEndian.Uint64(encodedNumberOfChunks.CandidateValue))

	entries := []node.DatastoreEntriesKeysAsString{}

	for i := 0; i < numberOfChunks; i++ {
		entry := node.DatastoreEntriesKeysAsString{
			Address: address,
			Key:     convert.EncodeStringToByteArray("massa_web_" + strconv.Itoa(i)),
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
