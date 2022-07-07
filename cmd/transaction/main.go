package main

import (
	"fmt"

	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/node/base58"
	sendOperation "github.com/massalabs/thyra/pkg/node/sendoperation"
	tx "github.com/massalabs/thyra/pkg/node/sendoperation/transaction"
)

func main() {
	base58Address := "A1pDYiLGT3Uzi8n1d2xzhuauzjb2BHZWbUMtQqKy7QeQGmEBLMW"

	addr, err := base58.CheckDecode(base58Address[1:])
	if err != nil {
		panic(err)
	}

	addr = addr[1:]

	pubKey, err := base58.CheckDecode("zkTGqfwJp43tY4FPgRXC7fr2xML3kDQ8bch15SpnDehuxWiKS")
	if err != nil {
		panic(err)
	}

	privKey, err := base58.CheckDecode("25CHWGN5DZemFnEdPyYfDkyYzEwierr3vCuP3Z4tiChfQpBP41")
	if err != nil {
		panic(err)
	}

	c := node.NewClient("http://localhost:33035") //https://test.massa.net/api/v2
	transaction := tx.New(addr, 100)

	id, err := sendOperation.Call(c, 31, 0, transaction, pubKey, privKey)
	if err != nil {
		panic(err)
	}

	fmt.Println("Execution OK, id is:", id)
}
