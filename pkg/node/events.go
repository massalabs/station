package node

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
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
	Slot              *Slot     `json:"slot"`
	Block             *string   `json:"block"`
	CallStack         *[]string `json:"call_stack"`
	ReadOnly          *bool     `json:"read_only"`
	IndexInSlot       *uint     `json:"index_in_slot"`
	OriginOperationID *string   `json:"origin_operation_id"`
}

const (
	maxWaitingTimeInSeconds = 45
	pollIntervalSec         = 1
)

/**
 * Filters events based on given arguments.
 *
 * Research criterion are
 * - by slots:
 * 	* after => start
 * 	* before => end
 * - callstack address:
 *  * last address in the callstack => emitter
 *  * first address in the callstack => originaCaller
 * - operation id.
 *
 * All these criterion are optional.
 */
func Events(client *Client, start *Slot, end *Slot,
	emitter *string, originaCaller *string,
	operationID *string,
) ([]Event, error) {
	rawResponse, err := client.RPCClient.Call(
		context.Background(),
		"get_filtered_sc_output_event",
		[]EventSearchCriteria{
			{
				Start:                 start,
				End:                   end,
				EmitterAddress:        emitter,
				OriginalCallerAddress: originaCaller,
				OriginalOperationID:   operationID,
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("calling endpoint get_filtered_sc_output_event with '%+v': %w",
			[]EventSearchCriteria{
				{
					Start:                 start,
					End:                   end,
					EmitterAddress:        emitter,
					OriginalCallerAddress: originaCaller,
					OriginalOperationID:   operationID,
				},
			}, err)
	}

	if rawResponse.Error != nil {
		return nil, rawResponse.Error
	}

	var resp []Event

	err = rawResponse.GetObject(&resp)
	if err != nil {
		return nil, fmt.Errorf("parsing get_filtered_sc_output_event jsonrpc response '%+v': %w", rawResponse, err)
	}

	return resp, nil
}

func ListenEvents(
	client *Client,
	start *Slot, end *Slot,
	emitter *string,
	operationID *string,
	caller *string,
) ([]Event, error) {
	counter := 0

	ticker := time.NewTicker(time.Second * pollIntervalSec)

	for ; true; <-ticker.C {
		counter++

		if counter > maxWaitingTimeInSeconds/pollIntervalSec {
			break
		}

		events, err := Events(client, start, end, emitter, caller, operationID)
		if err != nil {
			return nil,
				fmt.Errorf("fetching events for: opId %s, caller %s, emitter %s: %w", *operationID, *caller, *emitter, err)
		}

		for _, event := range events {
			if strings.Contains(event.Data, "massa_execution_error") {
				// return the event containing the error
				return nil, errors.New(event.Data)
			}
		}

		if len(events) > 0 {
			return events, nil
		}
	}

	return nil,
		fmt.Errorf("listening events for: opId %s, caller %s, emitter %s: Timeout", *operationID, *caller, *emitter)
}
