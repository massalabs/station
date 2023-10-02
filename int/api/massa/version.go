package massa

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station/api/swagger/server/models"
	"github.com/massalabs/station/api/swagger/server/restapi/operations"
)

func NewGetMassaStationVersion(version string) operations.GetMassaStationVersionHandler {
	return &getMassaStationVersion{version: version}
}

type getMassaStationVersion struct {
	version string
}

func (h *getMassaStationVersion) Handle(_ operations.GetMassaStationVersionParams) middleware.Responder {
	return operations.NewGetMassaStationVersionOK().
		WithPayload(
			models.Version(h.version),
		)
}
