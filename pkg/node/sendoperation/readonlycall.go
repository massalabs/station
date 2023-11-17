package sendoperation

import (
	"context"
	"fmt"

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

		// Don't return error, just return default gas limit,
		// because execute_read_only_call v27 execute_read_only_call may not work (Invalid params).
		return DefaultGasLimitCallSC, nil
	}

	estimatedGasCost := uint64(result.GasCost)

	return estimatedGasCost, nil
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
) (*ReadOnlyCallResult, error) {
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

	var resp []ReadOnlyCallResult

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
