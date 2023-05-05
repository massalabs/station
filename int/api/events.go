package api

import (
	"strings"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/models"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/config"
	"github.com/massalabs/thyra/pkg/node"
)

const (
	periodThreshold = 1
	tickerDuration  = time.Second
)

func NewEventListenerHandler(config *config.AppConfig) operations.ThyraEventsGetterHandler {
	return &eventListener{config: config}
}

type eventListener struct {
	config *config.AppConfig
}

//nolint:funlen
func (h *eventListener) Handle(params operations.ThyraEventsGetterParams) middleware.Responder {
	client := node.NewClient(h.config.NodeURL)

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

	ticker := time.NewTicker(tickerDuration)

	var event node.Event

	for ; true; <-ticker.C {
		trigger := false

		// In some cases, the event has been emitted in the previous period and we miss it.
		// So we check a given number of previous periods to be sure to catch it.
		slotStart.Period -= periodThreshold

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
