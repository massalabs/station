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
	secretKey := "S1pAYxBFbwomUUfzdv3yvQt8wPiYriX6VnXdrr87yrhbn91w96m"
	w, _ := wallet.NewFromSeed(secretKey)

	c := node.NewClient("https://test.massa.net/api/v2")
	//read smart contrat
	filePath := "path to your smart contract"
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	fmt.Print(data)

	exeSC := executesc.New(data, 700000, 2, 20)
	expirePeriod := uint64(37205) // TODO: obtain this from the network
	id, err := sendOperation.Call(c, expirePeriod, 0, exeSC, w.GetPublicKey(), w.GetPrivateKey())
	if err != nil {
		panic(err)
	}

	fmt.Println("Execution OK, id is:", id)
}
