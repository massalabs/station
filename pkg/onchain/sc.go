package onchain

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/massalabs/station/int/config"
	"github.com/massalabs/station/pkg/node"
	sendOperation "github.com/massalabs/station/pkg/node/sendoperation"
	"github.com/massalabs/station/pkg/node/sendoperation/callsc"
	"github.com/massalabs/station/pkg/node/sendoperation/executesc"
	"github.com/massalabs/station/pkg/node/sendoperation/signer"
)

type OperationWithEventResponse struct {
	Event             string
	OperationResponse sendOperation.OperationResponse
}

// CallFunction calls a function of a smart contract on the blockchain. It returns the operation ID or an Error if any.
// If `async` is true, it returns directly the operation ID and does not wait for the event.
// Otherwise, it waits for the first event generated by the smart contract and returns it along with the operation ID.
func CallFunction(
	networkInfos *config.NetworkInfos,
	nickname string,
	addr string,
	function string,
	parameter []byte,
	fee uint64,
	maxGas uint64,
	coins uint64,
	expiryDelta uint64,
	async bool,
	signer signer.Signer,
	description string,
) (*OperationWithEventResponse, error) {
	client := node.NewClient(networkInfos.NodeURL)

	// Calibrate max_gas
	if maxGas == 0 {
		estimatedGasCost, err := sendOperation.EstimateGasCostCallSC(nickname, addr, function, parameter, coins, fee, client)
		if err != nil {
			return nil, fmt.Errorf("estimating Call SC gas cost for function '%s' at '%s': %w", function, addr, err)
		}

		maxGas = estimatedGasCost
	}

	// Create the operation
	callSC, err := callsc.New(addr, function, parameter,
		maxGas,
		coins)
	if err != nil {
		return nil, fmt.Errorf("creating callSC with '%s' at '%s': %w", function, addr, err)
	}

	operationResponse, err := sendOperation.Call(
		client,
		networkInfos.ChainID,
		expiryDelta,
		fee,
		callSC,
		nickname,
		signer,
		description,
	)
	if err != nil {
		return nil, fmt.Errorf("calling function '%s' at '%s' with '%+v': %w", function, addr, parameter, err)
	}

	return CallFunctionSuccess(async, operationResponse, client)
}

func CallFunctionSuccess(
	async bool,
	operationResponse *sendOperation.OperationResponse,
	client *node.Client,
) (*OperationWithEventResponse, error) {
	if async {
		return &OperationWithEventResponse{
			Event:             "Function called successfully but did not wait for event",
			OperationResponse: *operationResponse,
		}, nil
	}

	events, err := node.ListenEvents(client, nil, nil, nil, &operationResponse.OperationID, nil, true)
	if err != nil {
		if strings.Contains(err.Error(), "Timeout") {
			return &OperationWithEventResponse{
				Event:             "Operation submited successfully but no event generated. The operation may have been rejected",
				OperationResponse: *operationResponse,
			}, nil
		}

		return nil, fmt.Errorf("listening events for opId at %s : %w", operationResponse.OperationID, err)
	}

	return &OperationWithEventResponse{
		Event:             events[0].Data, // return first event; TO DO: return all events
		OperationResponse: *operationResponse,
	}, nil
}

// DeploySC deploys a smart contract on the blockchain.
// The smart contract is deployed with the given account nickname.

func DeploySC(
	networkInfos *config.NetworkInfos,
	nickname string,
	maxGas uint64,
	maxCoins uint64,
	fees uint64,
	expiry uint64,
	parameters []byte,
	smartContractByteCode []byte,
	deployerByteCode []byte,
	signer signer.Signer,
	description string,
) (*sendOperation.OperationResponse, []node.Event, error) {
	client := node.NewClient(networkInfos.NodeURL)

	contract := DatastoreContract{
			Data:  smartContractByteCode,
			Args:  parameters,
			Coins: fees,
	}

	dataStore, err := populateDatastore(contract)
	if err != nil {
		return nil, nil, fmt.Errorf("populating datastore: %w", err)
	}


	exeSCOperation := executesc.New(
		deployerByteCode,
		maxGas,
		maxCoins,
		dataStore)
 		
	operationResponse, err := sendOperation.Call(
		client,
		networkInfos.ChainID,
		expiry,
		fees,
		exeSCOperation,
		nickname,
		signer,
		"Deploying smart contract: "+description,
	) 
	if err != nil {
		return nil, nil, fmt.Errorf("calling executeSC: %w", err)
	}

	events, err := node.ListenEvents(client, nil, nil, nil, &operationResponse.OperationID, nil, true)
	if err != nil {
		return nil, nil, fmt.Errorf("listening events for opId at %s : %w", operationResponse.OperationID, err)
	}

	return operationResponse, events, nil
}

func FindDeployedAddress(events []node.Event) (string, bool) {
	pattern := "Contract deployed at address: (.+)"
	re := regexp.MustCompile(pattern)

	for _, event := range events {
		match := re.FindStringSubmatch(event.Data)
		if len(match) > 1 {
			return match[1], true
		}
	}

	return "", false
}
