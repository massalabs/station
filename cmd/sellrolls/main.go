package main

import (
	"fmt"

	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/node/sendoperation"
	"github.com/massalabs/thyra/pkg/node/sendoperation/sellrolls"
	"github.com/massalabs/thyra/pkg/wallet"
)

func main() {

	w, _ := wallet.New("massa")
	c := node.NewClient()
	rolls := sellrolls.New(1)

	expirePeriod := uint64(36981)
	id, err := sendoperation.Call(c, expirePeriod, 0, rolls, w.KeyPairs[0].PrivateKey, w.KeyPairs[0].PrivateKey)
	if err != nil {
		panic(err)
	}
	fmt.Println("Execution OK, id is:", id)
}
