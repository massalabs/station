package sendoperation

import (
	"context"
	"fmt"
	"strconv"

	"github.com/massalabs/station/pkg/logger"
	"github.com/massalabs/station/pkg/node"
	"github.com/massalabs/station/pkg/wallet"
)

func EstimateGasCostCallSC(
	nickname string,
	targetAddr string,
	function string,
	parameter []byte,
	coins uint64,
	fee uint64,
	client *node.Client,
) (uint64, error) {
	acc, err := wallet.Fetch(nickname)
	if err != nil {
		return 0, fmt.Errorf("loading wallet '%s': %w", nickname, err)
	}

	result, err := ReadOnlyCallSC(targetAddr, function, parameter, coins, DefaultGasLimitCallSC, fee, acc.Address, client)
	if err != nil {
		logger.Warnf("calling ReadOnlyCall: %s", err)

		// Don't return error, just return default gas limit,
		// because execute_read_only_call v27 execute_read_only_call may not work (Invalid params)
		// and because execute_read_only_call v24.1 may fail because ignore the coins parameter.
		return DefaultGasLimitCallSC, nil
	}

	estimatedGasCost := uint64(result.GasCost)

	return estimatedGasCost, nil
}

func ReadOnlyCallSC(
	targetAddr string,
	function string,
	parameter []byte,
	coins uint64,
	maxGas uint64,
	fee uint64,
	callerAddr string,
	client *node.Client,
) (*ReadOnlyCallResult, error) {
	coinsString, err := NanoToMas(coins)
	if err != nil {
		return nil, fmt.Errorf("converting maxCoins to mas: %w", err)
	}

	feeString, err := NanoToMas(fee)
	if err != nil {
		return nil, fmt.Errorf("converting fee to mas: %w", err)
	}

	readOnlyCallParams := [][]ReadOnlyCallParams{
		{
			ReadOnlyCallParams{
				MaxGas:         int(maxGas),
				Coins:          coinsString,
				Fee:            feeString,
				TargetAddress:  targetAddr,
				TargetFunction: function,
				Parameter:      parameter,
				CallerAddress:  callerAddr,
			},
		},
	}

	rawResponse, err := client.RPCClient.Call(
		context.Background(),
		"execute_read_only_call",
		readOnlyCallParams,
	)
	if err != nil {
		return nil, fmt.Errorf("calling execute_read_only_call jsonrpc with '%+v': %w", readOnlyCallParams, err)
	}

	if rawResponse.Error != nil {
		return nil, fmt.Errorf(
			"receiving execute_read_only_call with '%+v': response: %w",
			readOnlyCallParams,
			rawResponse.Error,
		)
	}

	var resp []ReadOnlyCallResult

	err = rawResponse.GetObject(&resp)
	if err != nil {
		return nil, fmt.Errorf("parsing execute_read_only_call jsonrpc response '%+v': %w", rawResponse, err)
	}

	if resp[0].Result.Error != "" {
		return nil, fmt.Errorf("ReadOnlyCall error: %s, caller address is %s and coins are %d",
			resp[0].Result.Error, callerAddr, coins)
	}

	return &resp[0], nil
}

func EstimateGasCostExecuteSC(
	nickname string,
	contract []byte,
	datastore []byte,
	maxCoins uint64,
	fee uint64,
	client *node.Client,
) (uint64, error) {
	acc, err := wallet.Fetch(nickname)
	if err != nil {
		return 0, fmt.Errorf("loading wallet '%s': %w", nickname, err)
	}

	result, err := ReadOnlyExecuteSC(contract, datastore, maxCoins, MaxGasAllowedExecuteSC, fee, acc.Address, client)
	if err != nil {
		logger.Warnf("ReadOnlyExecuteSC error: %s, caller address is %s", err, acc.Address)

		// Do not return an error because execute_read_only_bytecode v24.1 fail with Internal error
		// and on v27 fails with Invalid params
		return DefaultGasLimitExecuteSC, nil
	}

	estimatedGasCost := uint64(result.GasCost)

	return estimatedGasCost, nil
}

func ReadOnlyExecuteSC(
	contract []byte,
	datastore []byte,
	maxCoins uint64,
	maxGas uint64,
	fee uint64,
	callerAddr string,
	client *node.Client,
) (*ReadOnlyCallResult, error) {
	if datastore == nil {
		datastore = []byte{0}
	}

	coins, err := NanoToMas(maxCoins)
	if err != nil {
		return nil, fmt.Errorf("converting maxCoins to mas: %w", err)
	}

	readOnlyExecuteParams := [][]ReadOnlyExecuteParams{
		{
			ReadOnlyExecuteParams{
				Coins:              coins,
				MaxGas:             int(maxGas),
				Fee:                strconv.FormatUint(fee, 10),
				Address:            callerAddr,
				Bytecode:           contract,
				OperationDatastore: datastore,
			},
		},
	}

	rawResponse, err := client.RPCClient.Call(
		context.Background(),
		"execute_read_only_bytecode",
		readOnlyExecuteParams,
	)
	if err != nil {
		return nil, fmt.Errorf("calling execute_read_only_bytecode jsonrpc: %w", err)
	}

	if rawResponse.Error != nil {
		return nil, fmt.Errorf(
			"receiving execute_read_only_bytecode  response: %w, %v",
			rawResponse.Error,
			fmt.Sprint(rawResponse.Error.Data),
		)
	}

	var resp []ReadOnlyCallResult // rename if it's the same for call sc and execute sc: ReadOnlyResult

	err = rawResponse.GetObject(&resp)
	if err != nil {
		return nil, fmt.Errorf("parsing execute_read_only_bytecode jsonrpc response '%+v': %w", rawResponse, err)
	}

	if resp[0].Result.Error != "" {
		return nil, fmt.Errorf("ReadOnlyExecuteSC error: %s, caller address is %s", resp[0].Result.Error, callerAddr)
	}

	return &resp[0], nil
}
