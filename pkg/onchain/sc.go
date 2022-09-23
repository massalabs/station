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

func CallFunction(client *node.Client, wallet wallet.Wallet,
	addr []byte, function string, parameter []byte,
) (string, error) {
	callSC := callsc.New(addr, function, parameter,
		sendOperation.NoGazFee, sendOperation.DefaultGazLimit,
		sendOperation.NoSequentialCoin, sendOperation.NoParallelCoin)

	operationID, err := sendOperation.Call(
		client,
		sendOperation.DefaultSlotsDuration, sendOperation.NoFee,
		callSC,
		wallet.KeyPairs[0].PublicKey, wallet.KeyPairs[0].PrivateKey)
	if err != nil {
		return "", fmt.Errorf("calling function '%s' at '%s' with '%+v': %w", function, addr, parameter, err)
	}

	counter := 0

	ticker := time.NewTicker(time.Minute)

	for ; true; <-ticker.C {
		counter++

		if counter > maxWaitingTimeInSeconds {
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

func CallFunctionUnwaited(client *node.Client, wallet wallet.Wallet,
	addr []byte, function string, parameter []byte,
) (string, error) {
	callSC := callsc.New(addr, function, parameter,
		sendOperation.NoGazFee, sendOperation.DefaultGazLimit,
		sendOperation.NoSequentialCoin, sendOperation.NoParallelCoin)

	operationID, err := sendOperation.Call(
		client,
		sendOperation.DefaultSlotsDuration, sendOperation.NoFee,
		callSC,
		wallet.KeyPairs[0].PublicKey, wallet.KeyPairs[0].PrivateKey)
	if err != nil {
		return "", fmt.Errorf("calling function '%s' at '%s' with '%+v': %w", function, addr, parameter, err)
	}

	return operationID, nil
}

func DeploySC(client *node.Client, wallet wallet.Wallet, contract []byte) (string, error) {
	exeSC := executesc.New(contract,
		sendOperation.DefaultGazLimit, sendOperation.NoGazFee,
		sendOperation.NoParallelCoin)

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

	ticker := time.NewTicker(time.Minute)

	for ; true; <-ticker.C {
		counter++

		if counter > maxWaitingTimeInSeconds {
			break
		}

		events, err := node.Events(client, nil, nil, nil, nil, &opID)
		if err != nil {
			return "", fmt.Errorf("waiting SC deployment: %w", err)
		}

		if len(events) > 0 {
			return strings.Split(events[0].Data, ":")[1], nil
		}
	}

	return "", errors.New("deployment time is out")
}
