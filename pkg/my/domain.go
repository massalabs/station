package my

import (
	"encoding/json"
	"fmt"

	"github.com/massalabs/thyra/api/swagger/server/models"
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

	domainsEntry, err := node.DatastoreEntry(client, dns.DNSRawAddress, ownedPrefix+wallet.Address)
	if err != nil {
		return nil, fmt.Errorf("reading entry '%s' at '%s': %w", dns.DNSRawAddress, ownedPrefix+wallet.Address, err)
	}

	err = json.Unmarshal(domainsEntry.CandidateValue, &domains)
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
			Key:     recordPrefix + domainNames[i],
		}
		params = append(params, param)
	}

	responses := []*models.Websites{}

	contractAddresses, err := node.DatastoreEntries(client, params)
	if err != nil {
		return nil, fmt.Errorf("reading entries'%s': %w", params, err)
	}

	for i := 0; i < len(domainNames); i++ {
		response := models.Websites{
			Address: string(contractAddresses[i].CandidateValue),
			Name:    domainNames[i],
		}
		responses = append(responses, &response)
	}

	return responses, nil
}
