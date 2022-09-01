package onchain

import (
	"errors"
	"strings"
	"time"

	"github.com/massalabs/thyra/pkg/node"
	sendOperation "github.com/massalabs/thyra/pkg/node/sendoperation"
	"github.com/massalabs/thyra/pkg/node/sendoperation/callsc"
	"github.com/massalabs/thyra/pkg/node/sendoperation/executesc"
	"github.com/massalabs/thyra/pkg/wallet"
)

func CallFunction(client *node.Client, wallet wallet.Wallet,
	addr []byte, function string, parameter []byte) (string, error) {
	callSC := callsc.New(addr, function, parameter, 0, 700000000, 0, 0)

	operationID, err := sendOperation.Call(client, 2, 0, callSC,
		wallet.KeyPairs[0].PublicKey,
		wallet.KeyPairs[0].PrivateKey)
	if err != nil {
		return "", err
	}

	counter := 0
	for range time.Tick(time.Second * 1) {
		counter++

		if counter > 45 {
			break
		}

		events, err := node.Events(client, nil, nil, nil, nil, &operationID)
		if err != nil {
			return operationID, err
		}

		if len(events) > 0 {
			return operationID, nil
		}
	}

	return operationID, errors.New("timeout")
}

func DeploySC(client *node.Client, wallet wallet.Wallet, contract []byte) (string, error) {
	exeSC := executesc.New([]byte(contract), 700000, 0, 0)

	id, err := sendOperation.Call(client, 2, 0, exeSC, wallet.KeyPairs[0].PublicKey, wallet.KeyPairs[0].PrivateKey)
	if err != nil {
		return "", err
	}

	counter := 0
	for range time.Tick(time.Second * 1) {
		counter++

		if counter > 45 {
			break
		}

		events, err := node.Events(client, nil, nil, nil, nil, &id)
		if err != nil {
			return "", err
		}

		if len(events) > 0 {
			return strings.Split(events[0].Data, ":")[1], nil
		}
	}

	return "", errors.New("deployment time is out.")
}
