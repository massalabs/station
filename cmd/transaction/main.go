package main

import (
	"fmt"

	"github.com/massalabs/thyra/pkg/node"
	sendOperation "github.com/massalabs/thyra/pkg/node/sendoperation"
	"github.com/massalabs/thyra/pkg/node/sendoperation/transaction"
	"github.com/massalabs/thyra/pkg/wallet"
)

func main() {
	wlt, err := wallet.New("massa")
	if err != nil {
		panic(err)
	}

	wlt2, err := wallet.New("massa")
	if err != nil {
		panic(err)
	}

	client := node.NewDefaultClient()

	trx := transaction.New(wlt2.KeyPairs[0].PublicKey, 1)

	opID, err := sendOperation.Call(
		client,
		sendOperation.DefaultSlotsDuration, sendOperation.NoFee,
		trx,
		wlt.KeyPairs[0].PublicKey, wlt.KeyPairs[0].PrivateKey)
	if err != nil {
		panic(err)
	}

	//nolint:forbidigo
	fmt.Println("Execution OK, id is:", opID)
}
