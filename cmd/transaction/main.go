package main

import (
	"fmt"

	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/node/sendoperation"
	"github.com/massalabs/thyra/pkg/node/sendoperation/transaction"
	"github.com/massalabs/thyra/pkg/wallet"
)

func main() {

	w, _ := wallet.New("massa")
	w2, _ := wallet.New("massa")
	c := node.NewClient("https://test.massa.net/api/v2")
	tx := transaction.New(w2.KeyPairs[0].PublicKey, 5)

	expirePeriod := uint64(37090)

	id, err := sendoperation.Call(c, expirePeriod, 0, tx, w.KeyPairs[0].PublicKey, w.KeyPairs[0].PrivateKey)
	if err != nil {
		panic(err)
	}

	fmt.Println("Execution OK, id is:", id)

}
