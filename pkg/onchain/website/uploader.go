package website

import (
	"encoding/json"
	"fmt"

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
	wallet, err := wallet.Load(walletNickname)
	if err != nil {
		return "", err
	}

	// Prepare address to webstorage.
	scAddress, err := onchain.DeploySC(client, wallet, []byte(sc.WebsiteStorer))
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	// Set DNS.
	_, err = dns.SetRecord(client, wallet, url, scAddress)
	if err != nil {
		fmt.Println(err)
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
	wallet, err := wallet.Load(walletNickname)
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

	return onchain.CallFunction(client, wallet, addr, "initializeWebsite", param)
}
