package main

import (
	"embed"
	"fmt"

	"github.com/massalabs/thyra/pkg/node"
	sendOperation "github.com/massalabs/thyra/pkg/node/sendoperation"
	"github.com/massalabs/thyra/pkg/node/sendoperation/executesc"
	"github.com/massalabs/thyra/pkg/wallet"
)

//go:embed sc
var content embed.FS

func main() {
	wlt, err := wallet.New("massa")
	if err != nil {
		panic(err)
	}

	client := node.NewDefaultClient()

	basePath := "sc/"

	websiteStorer, err := content.ReadFile(basePath + "websiteStorer.wasm")
	if err != nil {
		panic(err)
	}

	exeSC := executesc.New(websiteStorer,
		sendOperation.DefaultGazLimit,
		0, nil)

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
