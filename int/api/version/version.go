package version

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station/api/swagger/server/models"
	"github.com/massalabs/station/api/swagger/server/restapi/operations"
	"github.com/massalabs/station/int/config"
)

func Handle(_ operations.GetMassaStationVersionParams) middleware.Responder {
	return operations.NewGetMassaStationVersionOK().WithPayload(models.Version(config.Version))
}
