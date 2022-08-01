package getters

import (
	"context"

	"github.com/massalabs/thyra/pkg/node"
)

type getScOutputEvent struct {
	Start                 *Slot   `json:"start"`
	End                   *Slot   `json:"end"`
	EmitterAddress        *string `json:"emitter_address"`
	OriginalCallerAddress *string `json:"original_caller_address"`
	OriginalOperationId   *string `json:"original_operation_id"`
}

type getScOutputEventResponse struct {
	Data    string  `json:"data"`
	Context Context `json:"context"`
}

type Context struct {
	Slot                *Slot     `json:"slot"`
	Block               *string   `json:"block"`
	CallStack           *[]string `json:"call_stack"`
	ReadOnly            *bool     `json:"read_only"`
	IndexInSlot         *uint     `json:"index_in_slot"`
	OriginalOperationId *string   `json:"origin_operation_id"`
}
type Slot struct {
	Period int `json:"period"`
	Thread int `json:"thread"`
}

func GetEvents(client *node.Client, start *Slot, end *Slot, emiterAddress *string, callerAddress *string, operationId *string) (*[]getScOutputEventResponse, error) {
	response, err := client.RPCClient.Call(
		context.Background(),
		"get_filtered_sc_output_event",
		[]getScOutputEvent{
			{
				Start:                 start,
				End:                   end,
				EmitterAddress:        emiterAddress,
				OriginalCallerAddress: callerAddress,
				OriginalOperationId:   operationId,
			},
		},
	)

	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error
	}

	var entry *[]getScOutputEventResponse
	err = response.GetObject(&entry)
	if err != nil {
		return nil, err
	}

	return entry, nil
}
