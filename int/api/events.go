package api

import (
	"strings"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station/api/swagger/server/models"
	"github.com/massalabs/station/api/swagger/server/restapi/operations"
	"github.com/massalabs/station/int/config"
	"github.com/massalabs/station/pkg/node"
)

func NewEventListenerHandler(config *config.NetworkInfos) operations.EventsGetterHandler {
	return &eventListener{config: config}
}

type eventListener struct {
	config *config.NetworkInfos
}

const timeoutSec = 60

func (h *eventListener) Handle(params operations.EventsGetterParams) middleware.Responder {
	client := node.NewClient(h.config.NodeURL)

	status, err := node.Status(client)
	if err != nil {
		return operations.NewEventsGetterInternalServerError().WithPayload(
			&models.Error{
				Code:    errorCodeEventListener,
				Message: err.Error(),
			})
	}

	slotStart := node.Slot{
		// Listen to events from the last slot
		Period: status.LastSlot.Period - 1,
		Thread: 0,
	}

	start := time.Now()

	var event *node.Event

	for {
		events, err := node.ListenEvents(client, &slotStart, nil, nil, nil, &params.Caller, false)
		if err != nil {
			return operations.NewEventsGetterInternalServerError().
				WithPayload(
					&models.Error{
						Code:    errorCodeEventListener,
						Message: err.Error(),
					})
		}

		for index := range events {
			if strings.Contains(events[index].Data, params.Str) {
				event = &events[index]

				break
			}
		}

		if event != nil {
			break
		}

		elapsed := time.Since(start)
		if elapsed.Seconds() > timeoutSec {
			return operations.NewEventsGetterBadRequest().
				WithPayload(
					&models.Error{
						Code:    timeoutError,
						Message: "Event not found",
					})
		}
	}

	return operations.NewEventsGetterOK().WithPayload(&models.Events{
		Address: params.Caller,
		Data:    event.Data,
	})
}
