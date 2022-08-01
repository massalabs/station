package main

import (
	"fmt"
	"io/ioutil"

	"github.com/massalabs/thyra/pkg/node"
	sendOperation "github.com/massalabs/thyra/pkg/node/sendoperation"
	"github.com/massalabs/thyra/pkg/node/sendoperation/executesc"
	"github.com/massalabs/thyra/pkg/wallet"
)

func main() {

	w, _ := wallet.New("massa")
	c := node.NewClient()
	//read smart contrat
	filePath := "put here path to your smart contract"
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	exeSC := executesc.New(data, 700000, 2, 20)
	expirePeriod := uint64(52235) // TODO: obtain this from the network
	id, err := sendOperation.Call(c, expirePeriod, 0, exeSC, w.KeyPairs[0].PublicKey, w.KeyPairs[0].PrivateKey)
	if err != nil {
		panic(err)
	}

	fmt.Println("Execution OK, id is:", id)
}
