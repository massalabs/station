package main

import (
	"fmt"

	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/node/sendoperation"
	buyrolls "github.com/massalabs/thyra/pkg/node/sendoperation/buyrolls"
	"github.com/massalabs/thyra/pkg/wallet"
)

func main() {

	secretKey := "S1pAYxBFbwomUUfzdv3yvQt8wPiYriX6VnXdrr87yrhbn91w96m"
	w, _ := wallet.NewFromSeed(secretKey)
	c := node.NewClient("https://test.massa.net/api/v2")
	tx := buyrolls.New(1)

	id, err := sendoperation.Call(c, 36981, 0, tx, w.GetPublicKey(), w.GetPrivateKey())
	if err != nil {
		panic(err)
	}
	fmt.Println("Execution OK, id is:", id)
}
