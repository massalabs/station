package api

import (
	"strings"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/models"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/node"
)

func EventListenerHandler(params operations.ThyraEventsGetterParams) middleware.Responder {
	client := node.NewDefaultClient()

	status, err := node.Status(client)
	if err != nil {
		return operations.NewThyraEventsGetterInternalServerError().WithPayload(
			&models.Error{
				Code:    errorCodeEventListener,
				Message: err.Error(),
			})
	}

	slotStart := node.Slot{
		Period: status.LastSlot.Period,
		Thread: 0,
	}

	ticker := time.NewTicker(time.Second)

	var event node.Event

	for ; true; <-ticker.C {
		trigger := false

		events, err := node.Events(client, &slotStart, nil, &params.Caller, nil, nil)
		if err != nil {
			return operations.NewThyraEventsGetterInternalServerError().
				WithPayload(
					&models.Error{
						Code:    errorCodeEventListener,
						Message: err.Error(),
					})
		}

		for _, e := range events {
			if strings.Contains(e.Data, params.Str) {
				trigger = true
				event = e
			}
		}

		if trigger {
			break
		}

		status, err := node.Status(client)
		if err != nil {
			return operations.NewThyraEventsGetterInternalServerError().WithPayload(
				&models.Error{
					Code:    errorCodeEventListener,
					Message: err.Error(),
				})
		}

		slotStart.Period = status.LastSlot.Period
	}

	return operations.NewThyraEventsGetterOK().WithPayload(&models.Events{
		Address: params.Caller,
		Data:    event.Data,
	})
}
