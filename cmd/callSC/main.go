package main

import (
	"fmt"

	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/node/base58"
	"github.com/massalabs/thyra/pkg/node/send_operation"
	"github.com/massalabs/thyra/pkg/node/send_operation/call_sc"
)

func main() {
	base58Address := "A1MrqLgWq5XXDpTBH6fzXHUg7E8M5U2fYDAF3E1xnUSzyZuKpMh"
	addr_bytes, err := base58.CheckDecode(base58Address[1:])
	if err != nil {
		panic(err)
	}

	pubKey_bytes, err := base58.CheckDecode("zkTGqfwJp43tY4FPgRXC7fr2xML3kDQ8bch15SpnDehuxWiKS") //"zkTGqfwJp43tY4FPgRXC7fr2xML3kDQ8bch15SpnDehuxWiKS")
	if err != nil {
		panic(err)
	}

	privKey_bytes, err := base58.CheckDecode("25CHWGN5DZemFnEdPyYfDkyYzEwierr3vCuP3Z4tiChfQpBP41")
	if err != nil {
		panic(err)
	}
	c := node.NewClient("https://test.massa.net/api/v2")
	callSC := call_sc.New(addr_bytes, "set_dots", make([]byte, 0), 0, 700000000, 0, 0)
	id, err := send_operation.Call(c, 30903, 0, callSC, pubKey_bytes, privKey_bytes)
	if err != nil {
		panic(err)
	}
	fmt.Println("Execution OK, id is:", id)
}
