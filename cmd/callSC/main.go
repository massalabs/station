package main

import (
	"crypto/ed25519"
	"fmt"

	"github.com/massalabs/thyra/pkg/node"
	sendOperation "github.com/massalabs/thyra/pkg/node/sendoperation"
	callSC "github.com/massalabs/thyra/pkg/node/sendoperation/callsc"
)

func main() {
	//TODO https://github.com/massalabs/thyra/issues/27
	pubKey, privKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		panic(err)
	}

	c := node.NewClient()

	callSC := callSC.New(pubKey, "set_dots", make([]byte, 0), 0, 700000000, 0, 0)

	id, err := sendOperation.Call(c, 30903, 0, callSC, pubKey, privKey)
	if err != nil {
		panic(err)
	}

	fmt.Println("Execution OK, id is:", id)
}
