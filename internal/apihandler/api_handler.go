package apihandler

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/node/base58"
	sendOperation "github.com/massalabs/thyra/pkg/node/sendoperation"
	callSC "github.com/massalabs/thyra/pkg/node/sendoperation/callsc"
)

const errorCodeSendOperation = "0001"

func ExecuteFunction(params operations.CmdExecuteFunctionParams) middleware.Responder {
	addr, err := base58.CheckDecode((*params.Body.At)[1:])
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

	callSC := callSC.New(
		addr,
		*params.Body.Name,
		[]byte(params.Body.Args),
		uint64(params.Body.Gaz.Price),
		uint64(*params.Body.Gaz.Limit),
		uint64(params.Body.Coins.Sequential),
		uint64(params.Body.Coins.Parallel))

	c := node.NewClient("https://test.massa.net/api/v2")

	id, err := sendOperation.Call(
		c,
		30903,
		uint64(params.Body.Fee),
		callSC,
		pubKey, privKey)

	//errorCode := "toto"
	if err != nil {
		//operations.NewCmdExecuteFunctionInternalServerError().WithPayload(
		//	&operations.CmdExecuteFunctionInternalServerErrorBody{Code: &errorCode, Message: err.Error()})
		panic(err)

	}

	return operations.NewCmdExecuteFunctionOK().WithPayload(id)
}
