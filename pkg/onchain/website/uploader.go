package website

import (
	"encoding/json"

	"github.com/massalabs/thyra/pkg/my"
	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/node/base58"
	"github.com/massalabs/thyra/pkg/onchain"
	"github.com/massalabs/thyra/pkg/onchain/dns"
	"github.com/massalabs/thyra/pkg/sc"
	"github.com/massalabs/thyra/pkg/wallet"
)

func PrepareForUpload(url string, walletNickname string) (string, error) {
	client := node.NewDefaultClient()

	// TODO should use a nickname argument
	wallet, err := wallet.GetWallet(walletNickname)
	if err != nil {
		return "", err
	}

	// Prepare address to webstorage.
	scAddress, err := onchain.DeploySC(client, *wallet, []byte(sc.WebsiteStorer))
	if err != nil {
		return "", err
	}

	// Set DNS.
	_, err = dns.SetRecord(client, *wallet, url, scAddress)
	if err != nil {
		return "", err
	}

	// Add new record to my.Domains and save it for next time.
	dom, err := my.NewDomains()
	if err != nil {
		return "", err
	}

	dom.Add(my.Domain{URL: url, Address: scAddress})

	err = dom.Save()
	if err != nil {
		return "", err
	}

	return scAddress, nil
}

type UploadWebsiteParam struct {
	Data string `json:"data"`
}

func Upload(atAddress string, content string, walletNickname string) (string, error) {
	client := node.NewDefaultClient()

	// TODO should use a nickname argument
	wallet, err := wallet.GetWallet(walletNickname)
	if err != nil {
		return "", err
	}

	addr, _, err := base58.VersionedCheckDecode(atAddress[1:])
	if err != nil {
		return "", err
	}

	param, err := json.Marshal(UploadWebsiteParam{
		Data: content,
	})
	if err != nil {
		return "", err
	}

	return onchain.CallFunction(client, *wallet, addr, "initializeWebsite", param)
}
