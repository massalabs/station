package main

import (
	"fmt"

	"github.com/massalabs/thyra/pkg/node"
	sendOperation "github.com/massalabs/thyra/pkg/node/sendoperation"
	"github.com/massalabs/thyra/pkg/node/sendoperation/executesc"
	"github.com/massalabs/thyra/pkg/sc"
	"github.com/massalabs/thyra/pkg/wallet"
)

func main() {
	wlt, err := wallet.New("massa")
	if err != nil {
		panic(err)
	}

	client := node.NewDefaultClient()

	exeSC := executesc.New([]byte(sc.WebsiteStorer), 700000, 2, 20)

	id, err := sendOperation.Call(client, 2, 0, exeSC, wlt.KeyPairs[0].PublicKey, wlt.KeyPairs[0].PrivateKey)
	if err != nil {
		panic(err)
	}

	fmt.Println("Execution OK, id is:", id)
}
