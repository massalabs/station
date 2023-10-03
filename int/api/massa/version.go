package massa

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station/api/swagger/server/models"
	"github.com/massalabs/station/api/swagger/server/restapi/operations"
	"github.com/massalabs/station/int/config"
)

type GetMassaStationVersionFunc func(params operations.GetMassaStationVersionParams) middleware.Responder

func NewGetMassaStationVersion() GetMassaStationVersionFunc {
	return func(_ operations.GetMassaStationVersionParams) middleware.Responder {
		return operations.NewGetMassaStationVersionOK().
			WithPayload(models.Version(config.Version))
	}
}

func (f GetMassaStationVersionFunc) Handle(params operations.GetMassaStationVersionParams) middleware.Responder {
	return f(params)
}
