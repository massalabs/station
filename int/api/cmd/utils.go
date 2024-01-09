package cmd

import (
	"strconv"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station/api/swagger/server/models"
	"github.com/massalabs/station/api/swagger/server/restapi/operations"
	sendOperation "github.com/massalabs/station/pkg/node/sendoperation"
)

// amountToUint64 converts the given amount to uint64.
func amountToUint64(amount models.Amount, defaultValue uint64) (uint64, middleware.Responder) {
	result := defaultValue

	if string(amount) != "" {
		parsedAmount, err := strconv.ParseUint(string(amount), 10, 64)
		if err != nil {
			return 0, operations.NewCmdReadOnlyCallSCBadRequest().WithPayload(
				&models.Error{
					Code:    errorInvalidArgs,
					Message: "Error during amount conversion: " + err.Error(),
				})
		}

		result = parsedAmount
	}

	return result, nil
}

func amountToString(amount models.Amount, defaultValue uint64) (string, middleware.Responder) {
	amountUint64, errResponse := amountToUint64(amount, defaultValue)
	if errResponse != nil {
		return "", errResponse
	}

	amountString, err := sendOperation.NanoToMas(amountUint64)
	if err != nil {
		return "", operations.NewCmdReadOnlyCallSCBadRequest().WithPayload(
			&models.Error{
				Code:    errorInvalidArgs,
				Message: "Error during amount conversion: " + err.Error(),
			})
	}

	return amountString, nil
}

// CreateReadOnlyResult Converts an instance of sendOperation.ReadOnlyResult struct to models.ReadOnlyResult struct
//
//nolint:funlen
func CreateReadOnlyResult(result sendOperation.ReadOnlyResult) models.ReadOnlyResult {
	model := models.ReadOnlyResult{
		ExecutedAt: &models.Timestamp{
			Period: int64(result.ExecutedAt.Period), // Convert int to int64
			Thread: int64(result.ExecutedAt.Thread), // Convert int to int64
		},
		GasCost:      int64(result.GasCost), // Convert int to int64
		OutputEvents: []*models.Event{},
		Result: &models.Result{
			Ok:    result.Result.Ok,
			Error: result.Result.Error,
		},
		StateChanges: &models.StateChange{
			AsyncPoolChanges:             result.StateChanges.AsyncPoolChanges,
			ExecutedDenunciationsChanges: result.StateChanges.ExecutedDenunciationsChanges,
			ExecutedOpsChanges:           result.StateChanges.ExecutedOpsChanges,
			ExecutionTrailHashChange:     result.StateChanges.ExecutionTrailHashChange,
			LedgerChanges:                make(map[string]models.LedgerEntryChange),
			PosChanges: &models.PosChanges{
				DeferredCredits: &models.DeferredCreditsInfo{
					Credits: result.StateChanges.PosChanges.DeferredCredits.Credits,
				},
				ProductionStats: result.StateChanges.PosChanges.ProductionStats,
				RollChanges:     result.StateChanges.PosChanges.RollChanges,
				SeedBits: &models.SeedBitsInfo{
					Bits: int64(result.StateChanges.PosChanges.SeedBits.Bits), // Convert int to int64
					Data: result.StateChanges.PosChanges.SeedBits.Data,
					Head: &models.BitVecHeadInfo{
						Index: int64(result.StateChanges.PosChanges.SeedBits.Head.Index), // Convert int to int64
						Width: int64(result.StateChanges.PosChanges.SeedBits.Head.Width), // Convert int to int64
					},
					Order: result.StateChanges.PosChanges.SeedBits.Order,
				},
			},
		},
	}

	// Since OutputEvents and LedgerChanges are slices/maps of structs,
	// you need to loop through them and convert each individually.
	for _, oldEvent := range result.OutputEvents {
		newEvent := &models.Event{
			Context: &models.EventContext{
				Block:             oldEvent.Context.Block,
				CallStack:         oldEvent.Context.CallStack,
				IndexInSlot:       int64(oldEvent.Context.IndexInSlot), // Convert int to int64
				IsError:           oldEvent.Context.IsError,
				IsFinal:           oldEvent.Context.IsFinal,
				OriginOperationID: oldEvent.Context.OriginOperationID,
				ReadOnly:          oldEvent.Context.ReadOnly,
				Slot: &models.Timestamp{
					Period: int64(oldEvent.Context.Slot.Period), // Convert int to int64
					Thread: int64(oldEvent.Context.Slot.Thread), // Convert int to int64
				},
			},
			Data: oldEvent.Data,
		}
		model.OutputEvents = append(model.OutputEvents, newEvent)
	}

	for key, oldChange := range result.StateChanges.LedgerChanges {
		newChange := models.LedgerEntryChange{
			Update: &models.LedgerUpdate{
				Balance: &models.ChangeSet{
					Set: oldChange.Update.Balance.Set,
				},
				Bytecode:  oldChange.Update.Bytecode,
				Datastore: oldChange.Update.Datastore,
			},
		}
		model.StateChanges.LedgerChanges[key] = newChange
	}

	return model
}
