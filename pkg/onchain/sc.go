package onchain

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/massalabs/station/pkg/node"
	sendOperation "github.com/massalabs/station/pkg/node/sendoperation"
	"github.com/massalabs/station/pkg/node/sendoperation/callsc"
	"github.com/massalabs/station/pkg/node/sendoperation/executesc"
	"github.com/massalabs/station/pkg/node/sendoperation/signer"
)

const maxWaitingTimeInSeconds = 45

const evenHeartbeat = 2

type OperationWithEventResponse struct {
	Event             string
	OperationResponse sendOperation.OperationResponse
}

/*
This function send a callSC request to the node.
After a successful execution, listening to the events emitted by the operation
should be done on front end side by the consumer.
However in the current state of things the easiest way to unblock us is to
listen to these events in MassaStation and return them as a response.
Hence this function also listen for the first event emitted in an OP and returns it.
*/
func CallFunction(client *node.Client,
	nickname string,
	addr string,
	function string,
	parameter []byte,
	coins uint64,
	operationBatch sendOperation.OperationBatch,
	signer signer.Signer,
) (*OperationWithEventResponse, error) {
	callSC, err := callsc.New(addr, function, parameter,
		sendOperation.DefaultGasLimit,
		coins)
	if err != nil {
		return nil, fmt.Errorf("creating callSC with '%s' at '%s': %w", function, addr, err)
	}

	operationResponse, err := sendOperation.Call(
		client,
		sendOperation.DefaultSlotsDuration, sendOperation.NoFee,
		callSC,
		nickname,
		operationBatch,
		signer,
	)
	if err != nil {
		return nil, fmt.Errorf("calling function '%s' at '%s' with '%+v': %w", function, addr, parameter, err)
	}

	eventFound, operationWithEventResponse, err := listenEvents(client, operationResponse)
	if eventFound {
		return operationWithEventResponse, err
	}

	// If no event received, return a message to announce sc is called
	return &OperationWithEventResponse{
		Event:             "Function called successfully but no event generated",
		OperationResponse: *operationResponse,
	}, nil
}

func CallFunctionUnwaited(client *node.Client,
	nickname string,
	expiryDelta uint64,
	addr string,
	function string,
	parameter []byte,
	operationBatch sendOperation.OperationBatch,
	signer signer.Signer,
) (*sendOperation.OperationResponse, error) {
	callSC, err := callsc.New(addr, function, parameter,
		sendOperation.DefaultGasLimit,
		sendOperation.HundredMassa)
	if err != nil {
		return nil, fmt.Errorf("creating callSC with '%s' at '%s': %w", function, addr, err)
	}

	operationResponse, err := sendOperation.Call(
		client,
		expiryDelta, sendOperation.NoFee,
		callSC,
		nickname,
		operationBatch,
		signer,
	)
	if err != nil {
		return nil, fmt.Errorf("calling function '%s' at '%s' with '%+v': %w", function, addr, parameter, err)
	}

	return operationResponse, nil
}

// DeploySC deploys a smart contract on the blockchain. It returns the address of the smart contract and an Error.
// The smart contract is deployed with the given account nickname.
func DeploySC(client *node.Client,
	nickname string,
	gasLimit uint64,
	coins uint64,
	fee uint64,
	expiry uint64,
	contract []byte,
	datastore []byte,
	operationBatch sendOperation.OperationBatch,
	signer signer.Signer,
) (*OperationWithEventResponse, error) {
	exeSC := executesc.New(
		contract,
		gasLimit,
		coins,
		datastore)

	operationResponse, err := sendOperation.Call(
		client,
		expiry,
		fee,
		exeSC,
		nickname,
		operationBatch,
		signer,
	)
	if err != nil {
		return nil, fmt.Errorf("calling executeSC: %w", err)
	}

	eventFound, operationWithEventResponse, err := listenEvents(client, operationResponse)
	if eventFound {
		return operationWithEventResponse, err
	}

	// If no event received, return a message to announce sc is deployed
	return &OperationWithEventResponse{
		Event:             "sc deployed successfully but no event received",
		OperationResponse: *operationResponse,
	}, nil
}

func listenEvents(
	client *node.Client,
	operationResponse *sendOperation.OperationResponse,
) (bool, *OperationWithEventResponse, error) {
	counter := 0

	ticker := time.NewTicker(time.Second * evenHeartbeat)

	for ; true; <-ticker.C {
		counter++

		if counter > maxWaitingTimeInSeconds*evenHeartbeat {
			break
		}

		events, err := node.Events(client, nil, nil, nil, nil, &operationResponse.OperationID)
		if err != nil {
			return true, nil, fmt.Errorf("waiting SC deployment: %w", err)
		}

		if len(events) > 0 {
			event := events[0].Data

			// Catch Run Time Error and return it
			if strings.Contains(event, "massa_execution_error") {
				// return the event containing the error
				return true, nil, errors.New(event)
			}

			// if there is an event, return the first event
			return true, &OperationWithEventResponse{
				Event:             event,
				OperationResponse: *operationResponse,
			}, nil
		}
	}

	return false, nil, nil
}
