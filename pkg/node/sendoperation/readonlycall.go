package sendoperation

import (
	"context"
	"fmt"

	"github.com/massalabs/station/pkg/logger"
	"github.com/massalabs/station/pkg/node"
	"github.com/massalabs/station/pkg/wallet"
)

const (
	PercentageGasLimit = 10
)

func addXPercentage(x, y uint64) uint64 {
	return x + x*y/100
}

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

	coinsString, err := NanoToMas(coins)
	if err != nil {
		return 0, fmt.Errorf("converting maxCoins to mas: %w", err)
	}

	feeString, err := NanoToMas(fee)
	if err != nil {
		return 0, fmt.Errorf("converting fee to mas: %w", err)
	}

	result, err := ReadOnlyCallSC(targetAddr, function, parameter, coinsString, feeString, acc.Address, client)
	if err != nil {
		logger.Warnf("calling ReadOnlyCall: %s", err)

		return 0, err
	}

	estimatedGasCost := uint64(result.GasCost)
	addxpercentage := addXPercentage(estimatedGasCost, PercentageGasLimit)
	logger.Warn("estimatedGasCost: %+v\n\n\n", estimatedGasCost)
	logger.Warn("addXPercentage: %d\n\n\n", addxpercentage)
	return addXPercentage(estimatedGasCost, PercentageGasLimit), nil
}

// ReadOnlyCallSC calls execute_read_only_call jsonrpc method.
// coins and fee must be in MAS.
func ReadOnlyCallSC(
	targetAddr string,
	function string,
	parameter []byte,
	coins string,
	fee string,
	callerAddr string,
	client *node.Client,
) (*ReadOnlyResult, error) {
	readOnlyCallParams := [][]ReadOnlyCallParams{
		{
			ReadOnlyCallParams{
				MaxGas:         MaxGasAllowedCallSC,
				Coins:          coins,
				Fee:            fee,
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

	var resp []ReadOnlyResult

	err = rawResponse.GetObject(&resp)
	if err != nil {
		return nil, fmt.Errorf("parsing execute_read_only_call jsonrpc response '%+v': %w", rawResponse, err)
	}

	if resp[0].Result.Error != "" {
		return nil, fmt.Errorf("ReadOnlyCall error: %s, caller address is %s and coins are %s",
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

	coinsString, err := NanoToMas(maxCoins)
	if err != nil {
		return 0, fmt.Errorf("converting maxCoins to mas: %w", err)
	}

	feeString, err := NanoToMas(fee)
	if err != nil {
		return 0, fmt.Errorf("converting fee to mas: %w", err)
	}

	result, err := ReadOnlyExecuteSC(contract, datastore, coinsString, feeString, acc.Address, client)
	if err != nil {
		return 0, fmt.Errorf("ReadOnlyExecuteSC error: %w, caller address is %s", err, acc.Address)
	}

	estimatedGasCost := uint64(result.GasCost)

	return addXPercentage(estimatedGasCost, PercentageGasLimit), nil
}

// ReadOnlyExecuteSC calls execute_read_only_bytecode jsonrpc method.
// coins and fee must be in MAS.
func ReadOnlyExecuteSC(
	contract []byte,
	datastore []byte,
	coins string,
	fee string,
	callerAddr string,
	client *node.Client,
) (*ReadOnlyResult, error) {
	if datastore == nil {
		datastore = []byte{0}
	}

	readOnlyExecuteParams := [][]ReadOnlyExecuteParams{
		{
			ReadOnlyExecuteParams{
				Coins:              coins,
				MaxGas:             MaxGasAllowedExecuteSC,
				Fee:                fee,
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

	var resp []ReadOnlyResult

	err = rawResponse.GetObject(&resp)
	if err != nil {
		return nil, fmt.Errorf("parsing execute_read_only_bytecode jsonrpc response '%+v': %w", rawResponse, err)
	}

	if resp[0].Result.Error != "" {
		return nil, fmt.Errorf("ReadOnlyExecuteSC error: %s, caller address is %s", resp[0].Result.Error, callerAddr)
	}

	return &resp[0], nil
}
