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

	exeSC := executesc.New([]byte(sc.WebsiteStorer),
		sendOperation.DefaultGazLimit, sendOperation.NoGazFee,
		sendOperation.NoParallelCoin)

	opID, err := sendOperation.Call(
		client,
		sendOperation.DefaultSlotsDuration, sendOperation.NoFee,
		exeSC,
		wlt.KeyPairs[0].PublicKey, wlt.KeyPairs[0].PrivateKey)
	if err != nil {
		panic(err)
	}

	//nolint:forbidigo
	fmt.Println("Execution OK, id is:", opID)
}
