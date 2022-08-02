package website

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/node/getters"
	sendOperation "github.com/massalabs/thyra/pkg/node/sendoperation"
	"github.com/massalabs/thyra/pkg/node/sendoperation/executesc"
	"github.com/massalabs/thyra/pkg/wallet"
)

func CreateWebsiteDeployer() (*string, error) {
	rpcClient := node.NewClient()
	status, err := getters.GetNodeStatus(rpcClient)

	if err != nil {
		return nil, err
	}
	wallets, err := wallet.ReadWallets()
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadFile("../../resource/contracts/main-website-deployer.wasm")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	wallet := wallets[0]
	exeSC := executesc.New(data, 700000, 0, 0)
	expirePeriod := uint64(status.NextSlot.Period + 2)
	id, err := sendOperation.Call(rpcClient, expirePeriod, 0, exeSC, wallet.KeyPairs[0].PublicKey, wallet.KeyPairs[0].PrivateKey)
	if err != nil {
		return nil, err
	}
	fmt.Println("operation is i")
	return &id, nil
}

func RefreshDeployers() (deployers []string, e error) {

	wallets, err := wallet.ReadWallets()
	if err != nil {
		return nil, err
	}
	rpcClient := node.NewClient()
	events, err := getters.GetEvents(rpcClient, nil, nil, nil, &wallets[0].Address, nil)
	if err != nil {
		return nil, err
	}
	aa := *events
	for _, element := range aa {
		if strings.HasPrefix(element.Data, "Website Deployer is deployed at") {
			deployers = append(deployers, strings.Split(element.Data, ":")[1])
		}
	}

	bytesOutput, err := json.Marshal(deployers)
	if err != nil {
		return nil, err
	}

	err = os.WriteFile("deployers.json", bytesOutput, 0644)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	fmt.Println("asas")
	return nil, err
}

type Bytes struct {
	Data []byte `json:"data"`
}

// func UploadWebsite(contractAddress string) (s *string, err error) {

// 	c := node.NewClient()
// 	wallets, err := wallet.ReadWallets()
// 	if err != nil {
// 		return nil, err
// 	}
// 	status, err := getters.GetNodeStatus(c)
// 	if err != nil {
// 		return nil, err
// 	}

// 	address, _, err := base58.VersionedCheckDecode(contractAddress[1:])
// 	if err != nil {
// 		return nil, err
// 	}

// 	buf, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var entry Bytes

// 	entry.Data = buf
// 	data, err := json.Marshal(entry)
// 	if err != nil {
// 		return nil, err
// 	}

// 	callSC := callsc.New(address, "initializeWebsite", data, 0, 700000000, 0, 0)
// 	id, err := sendOperation.Call(c, uint64(status.NextSlot.Period+2), 0, callSC, wallets[0].KeyPairs[0].PublicKey, wallets[0].KeyPairs[0].PrivateKey)

// 	if err != nil {
// 		return nil, err
// 	}
// 	return &id, nil
// }
