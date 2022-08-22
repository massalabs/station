package dns

import (
	"encoding/json"
	"errors"

	"github.com/massalabs/thyra/api/swagger/server/models"
	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/node/base58"
	"github.com/massalabs/thyra/pkg/onchain"
	"github.com/massalabs/thyra/pkg/wallet"
)

const dnsRawAddress = "A12jkDPTcdhkqGg9VoKsTwvkBwZeSHQw7wJqQYKrNesKnjnGejuR"

func Resolve(client *node.Client, name string) (string, error) {
	const dnsPrefix = "record"

	entry, err := node.DatastoreEntry(client, dnsRawAddress, dnsPrefix+name)
	if err != nil {
		return "", err
	}
	if len(entry.CandidateValue) == 0 {
		return "", errors.New("name not found")
	}

	return string(entry.CandidateValue), nil
}

type setApproval struct {
	Operator string `json:"operator"`
	Approved bool   `json:"approved"`
}

type setRecord struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

func SetRecord(client *node.Client, wallet wallet.Wallet, url string, smartContract string) (string, error) {
	addr, _, err := base58.VersionedCheckDecode(dnsRawAddress[1:])
	if err != nil {
		return "", err
	}

	rec := setRecord{
		Name:    url,
		Address: smartContract,
	}

	param, err := json.Marshal(rec)
	if err != nil {
		return "", err
	}

	return onchain.CallFunction(client, wallet, addr, "setResolver", param)
}

func SetRecordManager(client *node.Client, wallet wallet.Wallet) (string, error) {
	addr, _, err := base58.VersionedCheckDecode(dnsRawAddress[1:])
	if err != nil {
		return "", err
	}

	// Set Resolver prepare data
	appr := &setApproval{
		Operator: wallet.Address,
		Approved: true,
	}

	param, err := json.Marshal(appr)
	if err != nil {
		return "", err
	}

	return onchain.CallFunction(client, wallet, addr, "setApprovalForAll", param)
}

func GetMyDomainNames(client *node.Client, nickname string) ([]string, error) {
	const ownedPrefix = "owned"
	wallet, err := wallet.GetWallet(nickname)
	if err != nil {
		return nil, err
	}
	domains := []string{}
	domainsEntry, err := node.DatastoreEntry(client, dnsRawAddress, ownedPrefix+wallet.Address)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(domainsEntry.CandidateValue, &domains)
	if err != nil {
		return nil, err
	}
	return domains, nil
}

func GetOwnedDomains(client *node.Client, domainNames []string) ([]*models.Websites, error) {
	const recordPrefix = "record"

	params := []node.GetDatastoreEntriesString{}
	for i := 0; i < len(domainNames); i++ {
		param := node.GetDatastoreEntriesString{
			Address: dnsRawAddress,
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
