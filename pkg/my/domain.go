package my

import (
	"encoding/json"

	"github.com/massalabs/thyra/api/swagger/server/models"
	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/onchain/dns"
	"github.com/massalabs/thyra/pkg/wallet"
)

const domainFile = "deployers.json"

type Domain struct {
	URL     string `json:"dnsName"`
	Address string `json:"address"`
}

type Domains struct {
	file    string
	domains []Domain
}

func GetDomains(client *node.Client, nickname string) ([]string, error) {
	const ownedPrefix = "owned"
	wallet, err := wallet.GetWallet(nickname)
	if err != nil {
		return nil, err
	}
	domains := []string{}
	domainsEntry, err := node.DatastoreEntry(client, dns.DnsRawAddress, ownedPrefix+wallet.Address)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(domainsEntry.CandidateValue, &domains)
	if err != nil {
		return nil, err
	}
	return domains, nil
}

func GetOwnedWebsites(client *node.Client, domainNames []string) ([]*models.Websites, error) {
	const recordPrefix = "record"

	params := []node.GetDatastoreEntriesString{}
	for i := 0; i < len(domainNames); i++ {
		param := node.GetDatastoreEntriesString{
			Address: dns.DnsRawAddress,
			Key:     recordPrefix + domainNames[i],
		}
		params = append(params, param)

	}

	responses := []*models.Websites{}
	contractAddresses, err := node.DatastoreEntries(client, params)

	contractAddressess := *contractAddresses
	for i := 0; i < len(domainNames); i++ {
		response := models.Websites{
			Address: string(contractAddressess[i].CandidateValue),
			Name:    domainNames[i],
		}
		responses = append(responses, &response)
	}
	if err != nil {
		return nil, err
	}

	return responses, nil
}
