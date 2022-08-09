package main

import (
	"fmt"

	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/node/sendoperation"
	"github.com/massalabs/thyra/pkg/node/sendoperation/buyrolls"
	"github.com/massalabs/thyra/pkg/wallet"
)

func main() {
	wlt, err := wallet.New("massa")
	if err != nil {
		panic(err)
	}

	client := node.NewDefaultClient()
	BuyRolls := buyrolls.New(1)

	id, err := sendoperation.Call(client, 2, 0, BuyRolls, wlt.KeyPairs[0].PublicKey, wlt.KeyPairs[0].PrivateKey)
	if err != nil {
		panic(err)
	}

	fmt.Println("Execution OK, id is:", id)
}
