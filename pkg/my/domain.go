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

func Domains(client *node.Client, nickname string) ([]string, error) {
	const ownedPrefix = "owned"

	wallet, err := wallet.Load(nickname)
	if err != nil {
		return nil, fmt.Errorf("loading wallet '%s': %w", nickname, err)
	}

	domains := []string{}

	keyWithHeader := convert.EncodeStringUint32ToUTF8(ownedPrefix + wallet.Address)

	domainsEntry, err := node.DatastoreEntry(client, dns.DNSRawAddress, keyWithHeader)
	if err != nil {
		return nil, fmt.Errorf("reading entry '%s' at '%s': %w", dns.DNSRawAddress, ownedPrefix+wallet.Address, err)
	}

	if len(domainsEntry.CandidateValue) == 0 {
		return domains, nil
	}

	domains = strings.Split(string(domainsEntry.CandidateValue[4:]), ",")

	if err != nil {
		return nil, fmt.Errorf("parsing json '%s': %w", domainsEntry.CandidateValue, err)
	}

	return domains, nil
}

func Websites(client *node.Client, domainNames []string) ([]*models.Websites, error) {
	const recordPrefix = "record"

	params := []node.DatastoreEntriesKeysAsString{}

	for i := 0; i < len(domainNames); i++ {
		param := node.DatastoreEntriesKeysAsString{
			Address: dns.DNSRawAddress,
			Key:     convert.EncodeStringUint32ToUTF8(recordPrefix + domainNames[i]),
		}
		params = append(params, param)
	}

	responses := []*models.Websites{}

	contractAddresses, err := node.DatastoreEntries(client, params)
	if err != nil {
		return nil, fmt.Errorf("reading entries'%s': %w", params, err)
	}

	for i := 0; i < len(domainNames); i++ { //nolint:varnamelen
		contractAddress := string(contractAddresses[i].CandidateValue[4:])

		brokenChunks, err := getMissingChunkIds(client, contractAddress)
		if err != nil {
			return nil, fmt.Errorf("checking chunk integrity: %w", err)
		}

		response := models.Websites{
			Address:      contractAddress,
			Name:         domainNames[i],
			BrokenChunks: brokenChunks,
		}
		responses = append(responses, &response)
	}

	return responses, nil
}

// Check website chunks and store its ID in an array if one of them is broken.
func getMissingChunkIds(client *node.Client, address string) ([]string, error) {
	chunkNumberKey := "total_chunks"

	var missedChunks []string

	keyNumber, err := node.DatastoreEntry(client, address, convert.EncodeStringUint32ToUTF8(chunkNumberKey))
	if err != nil {
		return nil, fmt.Errorf("reading datastore entry '%s' at '%s': %w", address, chunkNumberKey, err)
	}

	chunkNumber := int(binary.LittleEndian.Uint64(keyNumber.CandidateValue))

	entries := []node.DatastoreEntriesKeysAsString{}

	for i := 0; i < chunkNumber; i++ {
		entry := node.DatastoreEntriesKeysAsString{
			Address: address,
			Key:     convert.EncodeStringUint32ToUTF8("massa_web_" + strconv.Itoa(i)),
		}
		entries = append(entries, entry)
	}

	response, err := node.DatastoreEntries(client, entries)
	if err != nil {
		return nil, fmt.Errorf("calling get_datastore_entries '%+v': %w", entries, err)
	}

	for i := 0; i < chunkNumber; i++ {
		if string(response[i].CandidateValue) == "" {
			missedChunks = append(missedChunks, strconv.Itoa(i))
		}
	}

	return missedChunks, nil
}
