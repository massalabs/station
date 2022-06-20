package main

import (
	"fmt"

	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/node/base58"
	sendOperation "github.com/massalabs/thyra/pkg/node/sendoperation"
	callSC "github.com/massalabs/thyra/pkg/node/sendoperation/callsc"
)

func main() {
	base58Address := "A1MrqLgWq5XXDpTBH6fzXHUg7E8M5U2fYDAF3E1xnUSzyZuKpMh"

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

	c := node.NewClient("https://test.massa.net/api/v2")
	callSC := callSC.New(addr, "set_dots", make([]byte, 0), 0, 700000000, 0, 0)

	id, err := sendOperation.Call(c, 30903, 0, callSC, pubKey, privKey)
	if err != nil {
		panic(err)
	}

	fmt.Println("Execution OK, id is:", id)
}
