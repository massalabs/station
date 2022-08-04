package website

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/massalabs/thyra/pkg/contracts"
	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/node/base58"
	"github.com/massalabs/thyra/pkg/node/getters"
	sendOperation "github.com/massalabs/thyra/pkg/node/sendoperation"
	"github.com/massalabs/thyra/pkg/onchain/dns"

	"github.com/massalabs/thyra/pkg/node/sendoperation/callsc"
	"github.com/massalabs/thyra/pkg/node/sendoperation/executesc"
	"github.com/massalabs/thyra/pkg/wallet"
)

func deployWebsiteDeployer(c *node.Client, wallet wallet.Wallet, expire uint64) (*string, error) {
	exeSC := executesc.New([]byte(contracts.WebstiteDeployerContract), 700000, 0, 0)
	id, err := sendOperation.Call(c, expire, 0, exeSC, wallet.KeyPairs[0].PublicKey, wallet.KeyPairs[0].PrivateKey)
	if err != nil {
		return nil, err
	}

	// Get SC Contract
	smartContract := ""
	n := 0
	for n < 3 && smartContract == "" {

		time.Sleep(15 * time.Second)
		events, err := getters.GetEvents(c, nil, nil, nil, nil, &id)
		if err != nil {
			return nil, err
		}

		eventsValue := *events
		if len(eventsValue) > 0 {
			smartContract = strings.Split(eventsValue[0].Data, ":")[1]
		}
		n++

	}
	return &smartContract, nil

}

type websiteDeployer struct {
	DnsName *string `json:"dnsName"`
	Address *string `json:"address"`
}

func PostWebsite(dnsName string) (*string, error) {
	c := node.NewClient()

	// Get status for expire period
	expirePeriod, err := getters.GetExpirePeriod(c)
	if err != nil {
		return nil, err
	}

	// Get firt wallet
	wallet, err := wallet.GetFirstWallet()
	if err != nil {
		return nil, err
	}

	// Deploy Smart contract deployer
	smartContract, err := deployWebsiteDeployer(c, *wallet, *expirePeriod)
	if err != nil {
		return nil, err
	}

	// // Set DNS Approval
	// _, err = dns.SetDnsApproval(c, *wallet, true, *expirePeriod)
	// if err != nil {
	// 	return nil, err
	// }
	// time.Sleep(15 * time.Second)

	// Set DNS Resolver
	_, err = dns.SetDnsResolver(c, *wallet, dnsName, *smartContract, *expirePeriod)
	if err != nil {
		return nil, err
	}
	deployers, err := GetDeployers()
	if err != nil {
		return nil, err
	}
	dep := websiteDeployer{
		DnsName: &dnsName,
		Address: smartContract,
	}
	deployers = append(deployers, dep)
	fmt.Println(deployers)

	bytesOutput, err := json.Marshal(deployers)
	if err != nil {
		return nil, err
	}

	err = os.WriteFile("deployers.json", bytesOutput, 0o644)
	if err != nil {
		return nil, err
	}

	return smartContract, nil
}

func GetDeployers() ([]websiteDeployer, error) {
	bytesInput, err := ioutil.ReadFile("deployers.json")
	if err != nil {
		return nil, err
	}
	deployers := []websiteDeployer{}

	err = json.Unmarshal(bytesInput, &deployers)
	if err != nil {
		return nil, err
	}
	return deployers, nil
}

type Bytes struct {
	Data string `json:"data"`
}

func UploadWebsite(websiteData string, contractAddress string) (s *string, err error) {

	c := node.NewClient()
	wallets, err := wallet.ReadWallets()
	if err != nil {
		return nil, err
	}
	status, err := getters.GetNodeStatus(c)
	if err != nil {
		return nil, err
	}

	address, _, err := base58.VersionedCheckDecode(contractAddress[1:])
	if err != nil {
		return nil, err
	}

	d, err := json.Marshal(Bytes{
		Data: websiteData,
	})
	if err != nil {
		return nil, err
	}

	callSC := callsc.New(address, "initializeWebsite", d, 0, 700000000, 0, 0)
	id, err := sendOperation.Call(c, uint64(status.NextSlot.Period+2), 0, callSC, wallets[0].KeyPairs[0].PublicKey, wallets[0].KeyPairs[0].PrivateKey)

	if err != nil {
		return nil, err
	}
	b := false
	n := 0
	for n < 3 && !b {

		time.Sleep(15 * time.Second)
		events, err := getters.GetEvents(c, nil, nil, nil, nil, &id)
		if err != nil {
			return nil, err
		}

		eventsValue := *events
		if len(eventsValue) > 0 {
			b = true
		}
		n++

	}
	return &id, nil
}
