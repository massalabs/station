package massa

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station/api/swagger/server/models"
	"github.com/massalabs/station/api/swagger/server/restapi/operations"
	"github.com/massalabs/station/int/config"
)

func NewGetNodeHandler(config *config.NetworkInfos) operations.GetNodeHandler {
	return &getNodeHandler{config: config}
}

type getNodeHandler struct {
	config *config.NetworkInfos
}

func (h *getNodeHandler) Handle(_ operations.GetNodeParams) middleware.Responder {
	return operations.NewGetNodeOK().
		WithPayload(&models.MassaNodeItem{
			Network: h.config.Network,
			URL:     &h.config.NodeURL,
			ChainID: int64(h.config.ChainID),
		})
}
