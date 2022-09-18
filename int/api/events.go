package api

import (
	"strings"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/models"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/node"
)

//nolint:nolintlint,ireturn
func EventChecker(params operations.ThyraEventsGetterParams) middleware.Responder {
	const periodsBefore = 2

	client := node.NewDefaultClient()

	status, err := node.Status(client)
	if err != nil {
		return operations.NewThyraEventsGetterInternalServerError().
			WithPayload(
				&models.Error{
					Code:    errorCodeGetDomainNames,
					Message: err.Error(),
				})
	}

	slotStart := node.Slot{
		Period: status.LastSlot.Period - periodsBefore,
		Thread: 0,
	}

	ticker := time.NewTicker(time.Second)

	for ; true; <-ticker.C {
		trigger := false

		events, err := node.Events(client, &slotStart, nil, &params.Caller, nil, nil)
		if err != nil {
			return operations.NewThyraEventsGetterInternalServerError().
				WithPayload(
					&models.Error{
						Code:    errorCodeGetDomainNames,
						Message: err.Error(),
					})
		}

		for _, s := range events {
			if strings.Contains(s.Data, params.Str) {
				trigger = true
			}
		}

		if trigger {
			break
		}
	}

	return operations.NewThyraEventsGetterOK()
}
