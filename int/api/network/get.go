package network

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station/api/swagger/server/models"
	"github.com/massalabs/station/api/swagger/server/restapi/operations"
	"github.com/massalabs/station/int/config"
)

type getNetworkConfigHandler struct{ configManager *config.MSConfigManager }

func NewGetNetworkConfigHandler(configManager *config.MSConfigManager) operations.GetNetworkConfigHandler {
	return &getNetworkConfigHandler{configManager: configManager}
}

func (h *getNetworkConfigHandler) Handle(_ operations.GetNetworkConfigParams) middleware.Responder {
	currentNetwork := h.configManager.CurrentNetwork()
	availableNetworks := buildAvailableNetworkInfos(h.configManager)
	response := &models.NetworkManagerItem{
		CurrentNetwork:        &currentNetwork.Name,
		AvailableNetworkInfos: availableNetworks,
	}

	return operations.NewGetNetworkConfigOK().WithPayload(response)
}
