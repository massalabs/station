package main

import (
	"fmt"

	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/node/sendoperation"
	"github.com/massalabs/thyra/pkg/node/sendoperation/transaction"
	"github.com/massalabs/thyra/pkg/wallet"
)

func main() {

	secretKey := "S1pAYxBFbwomUUfzdv3yvQt8wPiYriX6VnXdrr87yrhbn91w96m"
	w, _ := wallet.NewFromSeed(secretKey)
	//w.Print()
	secretKey2 := "S1rRCHbzV3P7SYtTzQDzTasK3PXR94yyA5pvyxcYa4VFj4785AZ"
	w2, _ := wallet.NewFromSeed(secretKey2)
	//w2.Print()

	c := node.NewClient("https://test.massa.net/api/v2")
	tx := transaction.New(w2.Address, 5)

	expirePeriod := uint64(37090)

	id, err := sendoperation.Call(c, expirePeriod, 0, tx, w.GetPublicKey(), w.GetPrivateKey())
	if err != nil {
		panic(err)
	}

	fmt.Println("Execution OK, id is:", id)

}
