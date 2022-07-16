package main

import (
	"fmt"

	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/node/sendoperation"
	"github.com/massalabs/thyra/pkg/node/sendoperation/sellrolls"
	"github.com/massalabs/thyra/pkg/wallet"
)

func main() {

	secretKey := "S1pAYxBFbwomUUfzdv3yvQt8wPiYriX6VnXdrr87yrhbn91w96m" //make sure you use the key with coins
	w, _ := wallet.NewFromSeed(secretKey)
	c := node.NewClient("https://test.massa.net/api/v2")
	rolls := sellrolls.New(1)

	expirePeriod := uint64(36981)
	id, err := sendoperation.Call(c, expirePeriod, 0, rolls, w.GetPublicKey(), w.GetPrivateKey())
	if err != nil {
		panic(err)
	}
	fmt.Println("Execution OK, id is:", id)
}
