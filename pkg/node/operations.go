package node

import (
	"context"
	"fmt"

	"github.com/massalabs/station/pkg/node/sendoperation/buyrolls"
	"github.com/massalabs/station/pkg/node/sendoperation/callsc"
	"github.com/massalabs/station/pkg/node/sendoperation/executesc"
	"github.com/massalabs/station/pkg/node/sendoperation/sellrolls"
	"github.com/massalabs/station/pkg/node/sendoperation/transaction"
)

//nolint:tagliatelle
type Operation struct {
	ID       *string  `json:"id"`
	InBlocks []string `json:"in_blocks"`
	InPool   bool     `json:"in_pool"`
	IsFinal  bool     `json:"is_final"`
	Detail   *Detail  `json:"operation"`
}

type Detail struct {
	Content   Content `json:"content"`
	Signature string  `json:"signature"`
}

type Content struct {
	ExpirePeriod uint   `json:"expire_period"`
	Fee          string `json:"fee"`
	Op           Op     `json:"op"`
}

//nolint:tagliatelle
type Op struct {
	Transaction *transaction.Transaction `json:"Transaction"`
	RollBuy     *buyrolls.BuyRolls       `json:"RollBuy"`
	RollSell    *sellrolls.SellRolls     `json:"RollSell"`
	ExecuteSC   *executesc.ExecuteSC     `json:"ExecuteSC"`
	CallSC      *callsc.CallSC           `json:"CallSC"`
}

func Operations(client *Client, ids []string) ([]Operation, error) {
	rawResponse, err := client.RPCClient.Call(
		context.Background(),
		"get_operations",
		append([][]string{}, ids),
	)
	if err != nil {
		return nil, fmt.Errorf("calling get_operations with '%+v': %w", append([][]string{}, ids), err)
	}

	if rawResponse.Error != nil {
		return nil, rawResponse.Error
	}

	var resp []Operation

	err = rawResponse.GetObject(&resp)
	if err != nil {
		return nil, fmt.Errorf("parsing get_operations jsonrpc response '%+v': %w", rawResponse, err)
	}

	return resp, nil
}
