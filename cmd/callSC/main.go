package main

import (
	"fmt"

	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/node/base58"
	sendOperation "github.com/massalabs/thyra/pkg/node/sendoperation"
	callSC "github.com/massalabs/thyra/pkg/node/sendoperation/callsc"
	"github.com/massalabs/thyra/pkg/wallet"
)

func main() {
	rawAddr := "A1JEEbgWPQMt97pJmZ3akxU64yW82wmZRe8EefjnEPxyCHgv1Yn"

	wlt, err := wallet.New("massa")
	if err != nil {
		panic(err)
	}

	client := node.NewDefaultClient()

	addr, _, err := base58.VersionedCheckDecode(rawAddr[1:])
	if err != nil {
		panic(err)
	}

	callSC := callSC.New(addr, "set_dots", make([]byte, 0), 0, 700000000, 0, 0)

	id, err := sendOperation.Call(client, 2, 0, callSC, wlt.KeyPairs[0].PublicKey, wlt.KeyPairs[0].PrivateKey)
	if err != nil {
		panic(err)
	}

	fmt.Println("Execution OK, id is:", id)
}
