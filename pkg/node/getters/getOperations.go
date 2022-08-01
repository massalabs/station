package getters

import (
	"context"

	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/node/sendoperation/buyrolls"
	"github.com/massalabs/thyra/pkg/node/sendoperation/callsc"
	"github.com/massalabs/thyra/pkg/node/sendoperation/sellrolls"
	"github.com/massalabs/thyra/pkg/node/sendoperation/transaction"
)

type getOperations struct {
	Id        *string    `json:"id"`
	InBlocks  *[]string  `json:"in_blocks"`
	InPool    bool       `json:"in_pool"`
	IsFinal   bool       `json:"is_final"`
	Operation *Operation `json:"operation"`
}
type Operation struct {
	Content   Content `json:"content"`
	Signature string  `json:"signature"`
}

type Content struct {
	ExpirePeriod uint   `json:"expire_period"`
	Fee          string `json:"fee"`
	Op           Op     `json:"op"`
}

type Op struct {
	Transaction *transaction.Transaction `json:"Transaction"`
	RollBuy     *buyrolls.BuyRolls       `json:"RollBuy"`
	RollSell    *sellrolls.SellRolls     `json:"RollSell"`
	ExecuteSC   *ExecuteSC               `json:"ExecuteSC"`
	CallSC      *callsc.CallSC           `json:"CallSC"`
}

type ExecuteSC struct {
	data     []byte
	maxGas   uint64
	gasPrice uint64
	Coins    string
}

func GetOperations(client *node.Client, operations []string) (*[]getOperations, error) {
	var stringg [][]string
	stringg = append(stringg, operations)
	response, err := client.RPCClient.Call(
		context.Background(),
		"get_operations",
		stringg,
	)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error
	}
	var entry *[]getOperations
	err = response.GetObject(&entry)
	if err != nil {
		return nil, err
	}
	return entry, nil
}
