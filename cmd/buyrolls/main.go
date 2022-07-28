package main

import (
	"fmt"

	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/node/sendoperation"
	"github.com/massalabs/thyra/pkg/node/sendoperation/buyrolls"
	"github.com/massalabs/thyra/pkg/wallet"
)

func main() {

	w, _ := wallet.New("massa")
	c := node.NewClient("https://test.massa.net/api/v2")
	tx := buyrolls.New(1)

	id, err := sendoperation.Call(c, 36981, 0, tx, w.KeyPairs[0].PublicKey, w.KeyPairs[0].PrivateKey)
	if err != nil {
		panic(err)
	}
	fmt.Println("Execution OK, id is:", id)
}
