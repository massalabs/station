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
	"github.com/massalabs/thyra/pkg/wallet"
)

const maxWaitingTimeInSeconds = 45

const evenHeartbeat = 2

func CallFunction(client *node.Client, wallet wallet.Wallet,
	addr []byte, function string, parameter []byte, coins uint64,
) (string, error) {
	callSC := callsc.New(addr, function, parameter,
		sendOperation.DefaultGazLimit,
		coins)

	operationID, err := sendOperation.Call(
		client,
		sendOperation.DefaultSlotsDuration, sendOperation.NoFee,
		callSC,
		wallet.KeyPairs[0].PublicKey, wallet.KeyPairs[0].PrivateKey)
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
			return operationID,
				fmt.Errorf("waiting execution of function '%s' at '%s' with id '%s': %w", function, addr, operationID, err)
		}

		if len(events) > 0 {
			return operationID, nil
		}
	}

	return operationID, errors.New("timeout")
}

func CallFunctionUnwaited(client *node.Client, wallet wallet.Wallet, expiryDelta uint64,
	addr []byte, function string, parameter []byte,
) (string, error) {
	callSC := callsc.New(addr, function, parameter,
		sendOperation.DefaultGazLimit,
		sendOperation.HundredMassa)

	operationID, err := sendOperation.Call(
		client,
		expiryDelta, sendOperation.NoFee,
		callSC,
		wallet.KeyPairs[0].PublicKey, wallet.KeyPairs[0].PrivateKey)
	if err != nil {
		return "", fmt.Errorf("calling function '%s' at '%s' with '%+v': %w", function, addr, parameter, err)
	}

	return operationID, nil
}

func DeploySC(client *node.Client, wallet wallet.Wallet, contract []byte) (string, error) {
	datastore := make(map[[3]uint8][]uint8)

	datastore[[3]uint8{1, 2, 3}] = []uint8{1, 2, 3}
	exeSC := executesc.New(contract,
		sendOperation.DefaultGazLimit,
		sendOperation.NoCoin, datastore)

	opID, err := sendOperation.Call(
		client,
		sendOperation.DefaultSlotsDuration,
		sendOperation.NoFee,
		exeSC,
		wallet.KeyPairs[0].PublicKey, wallet.KeyPairs[0].PrivateKey)
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

// DeploySC deploys a smart contract on the blockchain. It returns the address of the smart contract and an Error.
// The smart contract is deployed with the given wallet.
func DeploySCV2(client *node.Client, nickname string, contract []byte) (string, error) {
	datastore := make(map[[3]uint8][]uint8)

	datastore[[3]uint8{1, 2, 3}] = []uint8{1, 2, 3}
	exeSC := executesc.New(contract,
		sendOperation.DefaultGazLimit,
		sendOperation.NoCoin, datastore)

	opID, err := sendOperation.CallV2(
		client,
		sendOperation.DefaultSlotsDuration,
		sendOperation.NoFee,
		exeSC,
		nickname)
	if err != nil {
		return "", fmt.Errorf("calling executeSC: %w", err)
	}

	counter := 0

	ticker := time.NewTicker(time.Second * evenHeartbeat)
	fmt.Println("Waiting for events...")
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
