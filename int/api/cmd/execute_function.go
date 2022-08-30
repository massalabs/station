package cmd

import (
	"sync"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/models"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/node/base58"
	sendOperation "github.com/massalabs/thyra/pkg/node/sendoperation"
	callSC "github.com/massalabs/thyra/pkg/node/sendoperation/callsc"
	"github.com/massalabs/thyra/pkg/wallet"
)

func NewExecuteFunction(walletStorage *sync.Map) operations.CmdExecuteFunctionHandler {
	return &functionExecuter{walletStorage: walletStorage}
}

type functionExecuter struct {
	walletStorage *sync.Map
}

//TODO Manage the panic(error)
func (f *functionExecuter) Handle(params operations.CmdExecuteFunctionParams) middleware.Responder {
	addr, err := base58.CheckDecode((*params.Body.At)[1:])
	if err != nil {
		panic(err)
	}

	addr = addr[1:]

	value, ok := f.walletStorage.Load(*params.Body.Name)
	if !ok {
		return operations.NewCmdExecuteFunctionUnprocessableEntity().WithPayload(
			&models.Error{
				Code:    errorCodeUnknownKeyID,
				Message: "Error: unknown wallet nickname : " + *params.Body.Name,
			})
	}

	storedWallet, ok := value.(*wallet.Wallet)
	if !ok {
		panic("stored value is not a wallet")
	}

	//TODO: storedWallet.KeyPairs[0].Protected shall be false
	pubKey := storedWallet.KeyPairs[0].PublicKey
	privKey := storedWallet.KeyPairs[0].PrivateKey

	callSC := callSC.New(
		addr,
		*params.Body.Name,
		[]byte(params.Body.Args),
		uint64(params.Body.Gaz.Price),
		uint64(*params.Body.Gaz.Limit),
		uint64(params.Body.Coins.Sequential),
		uint64(params.Body.Coins.Parallel))

	c := node.NewDefaultClient()

	operationID, err := sendOperation.Call(
		c,
		30903,
		uint64(params.Body.Fee),
		callSC,
		pubKey, privKey)
	if err != nil {
		return operations.NewCmdExecuteFunctionInternalServerError().WithPayload(
			&models.Error{Code: errorCodeSendOperation, Message: "Error: " + err.Error()})
	}

	return operations.NewCmdExecuteFunctionOK().WithPayload(operationID)
}
