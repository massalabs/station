package massa

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station/api/swagger/server/models"
	"github.com/massalabs/station/api/swagger/server/restapi/operations"
	"github.com/massalabs/station/int/config"
)

func GetMassaStationVersionFunc(_ operations.GetMassaStationVersionParams) middleware.Responder {
	return operations.NewGetMassaStationVersionOK().WithPayload(models.Version(config.Version))
}
