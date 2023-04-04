package onchain

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/massalabs/thyra/pkg/node"
	sendOperation "github.com/massalabs/thyra/pkg/node/sendoperation"
	"github.com/massalabs/thyra/pkg/node/sendoperation/callsc"
	"github.com/massalabs/thyra/pkg/node/sendoperation/executesc"
)

const maxWaitingTimeInSeconds = 45

const evenHeartbeat = 2

/*
This function send a callSC request to the node.
After a sucessfulll execution, listening to the events emitted by the operation
should be done on front end side by the consumer.
However in the current state of things the easiest way to unblock us is to
listen to these events in Thyra and return them as a response.
Hence this function also listen for the first event emitted in an OP and returns it.
*/
func CallFunction(client *node.Client, nickname string,
	addr []byte, function string, parameter []byte, coins uint64,
) (string, error) {
	callSC := callsc.New(addr, function, parameter,
		sendOperation.DefaultGazLimit,
		coins)

	operationID, err := sendOperation.Call(
		client,
		sendOperation.DefaultSlotsDuration, sendOperation.NoFee,
		callSC,
		nickname)
	if err != nil {
		return "", fmt.Errorf("calling function '%s' at '%s' with '%+v': %w", function, addr, parameter, err)
	}

	counter := 0

	ticker := time.NewTicker(time.Second * evenHeartbeat)

	for ; true; <-ticker.C {
		counter++

		if counter > maxWaitingTimeInSeconds*evenHeartbeat {
			break
		}

		events, err := node.Events(client, nil, nil, nil, nil, &operationID)
		if err != nil {
			return "", fmt.Errorf("listening events: %w", err)
		}

		if len(events) > 0 {
			event := events[0].Data
			//  Catch Run Time Error and return it
			if strings.Contains(event, "massa_execution_error") {
				// return the event containing the error
				return "", errors.New(event)
			}
			// if there is an event, return the first event
			return event, nil
		}
	}
	// If no event received, return a message to announce sc is deployed
	return "Function called successfully but no event generated", nil
}

func CallFunctionUnwaited(client *node.Client, nickname string, expiryDelta uint64,
	addr []byte, function string, parameter []byte,
) (string, error) {
	callSC := callsc.New(addr, function, parameter,
		sendOperation.DefaultGazLimit,
		sendOperation.HundredMassa)

	operationID, err := sendOperation.Call(
		client,
		expiryDelta, sendOperation.NoFee,
		callSC,
		nickname)
	if err != nil {
		return "", fmt.Errorf("calling function '%s' at '%s' with '%+v': %w", function, addr, parameter, err)
	}

	return operationID, nil
}

// DeploySC deploys a smart contract on the blockchain. It returns the address of the smart contract and an Error.
// The smart contract is deployed with the given account nickname.
func DeploySC(client *node.Client,
	nickname string,
	gazLimit uint64,
	coins uint64,
	fee uint64,
	expiry uint64,
	contract []byte,
	datastore []byte,
) (string, error) {
	exeSC := executesc.New(
		contract,
		gazLimit,
		coins,
		datastore)

	opID, err := sendOperation.Call(
		client,
		expiry,
		fee,
		exeSC,
		nickname)
	if err != nil {
		return "", fmt.Errorf("calling executeSC: %w", err)
	}

	counter := 0

	ticker := time.NewTicker(time.Second * evenHeartbeat)

	for ; true; <-ticker.C {
		counter++

		if counter > maxWaitingTimeInSeconds*evenHeartbeat {
			break
		}

		events, err := node.Events(client, nil, nil, nil, nil, &opID)
		if err != nil {
			return "", fmt.Errorf("waiting SC deployment: %w", err)
		}

		if len(events) > 0 {
			event := events[0].Data
			//  Catch Run Time Error and return it
			if strings.Contains(event, "massa_execution_error") {
				// return the event containing the error
				return "", errors.New(event)
			}
			// if there is an event, return the first event
			return event, nil
		}
	}
	// If no event received, return a message to announce sc is deployed
	return "sc deployed successfully but no event received", nil
}
