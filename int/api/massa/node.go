package massa

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station/api/swagger/server/models"
	"github.com/massalabs/station/api/swagger/server/restapi/operations"
	"github.com/massalabs/station/int/config"
)

func NewGetNodeHandler(configManager *config.MSConfigManager) operations.GetNodeHandler {
	return &getNodeHandler{configManager: configManager}
}

type getNodeHandler struct{ configManager *config.MSConfigManager }

func (h *getNodeHandler) Handle(_ operations.GetNodeParams) middleware.Responder {
	currentNetwork := h.configManager.CurrentNetwork()
	return operations.NewGetNodeOK().
		WithPayload(&models.MassaNodeItem{
			Network: &currentNetwork.Name,
			URL:     &currentNetwork.NodeURL,
			ChainID: int64(currentNetwork.ChainID),
		})
}
