package node

import (
	"context"
)

type EventSearchCriteria struct {
	Start                 *Slot   `json:"start"`
	End                   *Slot   `json:"end"`
	EmitterAddress        *string `json:"emitter_address"`
	OriginalCallerAddress *string `json:"original_caller_address"`
	OriginalOperationID   *string `json:"original_operation_id"`
}

type Event struct {
	Data    string  `json:"data"`
	Context Context `json:"context"`
}

type Context struct {
	Slot                *Slot     `json:"slot"`
	Block               *string   `json:"block"`
	CallStack           *[]string `json:"call_stack"`
	ReadOnly            *bool     `json:"read_only"`
	IndexInSlot         *uint     `json:"index_in_slot"`
	OriginalOperationID *string   `json:"origin_operation_id"`
}

/**
 * Filters events based on given arguments.
 *
 * Research criterion are
 * - by slots:
 * 	* after => start
 * 	* before => end
 * - callstack address:
 *  * last address in the callstack => trigger
 *  * first address in the callstack => originator
 * - operation id.
 *
 * All these criterion are optional.
 */
func Events(client *Client, start *Slot, end *Slot,
	trigger *string, originator *string,
	operationID *string,
) ([]Event, error) {
	rawResponse, err := client.RPCClient.Call(
		context.Background(),
		"get_filtered_sc_output_event",
		[]EventSearchCriteria{
			{
				Start:                 start,
				End:                   end,
				EmitterAddress:        trigger,
				OriginalCallerAddress: originator,
				OriginalOperationID:   operationID,
			},
		},
	)
	if err != nil {
		return nil, err
	}

	if rawResponse.Error != nil {
		return nil, rawResponse.Error
	}

	var resp []Event
	err = rawResponse.GetObject(&resp)

	if err != nil {
		return nil, err
	}

	return resp, nil
}
