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
	eventPollingTimeoutSec = 60
	pollIntervalSec        = 1
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
	emitter *string, originalCaller *string,
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
				OriginalCallerAddress: originalCaller,
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
					OriginalCallerAddress: originalCaller,
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
	failOnExecError bool,
) ([]Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(eventPollingTimeoutSec))
	defer cancel()

	ticker := time.NewTicker(time.Second * pollIntervalSec)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("no events found in the given time interval (%d seconds)", eventPollingTimeoutSec)
		case <-ticker.C:
			events, err := Events(client, start, end, emitter, caller, operationID)
			if err != nil {
				return nil, fmt.Errorf("fetching events: %w", err)
			}

			for _, event := range events {
				if failOnExecError && strings.Contains(event.Data, "massa_execution_error") {
					// return the event containing the error
					return nil, errors.New(event.Data)
				}
			}

			if len(events) > 0 {
				return events, nil
			}
		}
	}
}
