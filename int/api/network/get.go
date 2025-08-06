package network

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station/api/swagger/server/models"
	"github.com/massalabs/station/api/swagger/server/restapi/operations"
	"github.com/massalabs/station/int/config"
)

type getNetworkConfigHandler struct{ configManager *config.MSConfigManager }

// NewGetNetworkConfigHandler creates a new getNetworkConfigHandler instance.
func NewGetNetworkConfigHandler(configManager *config.MSConfigManager) operations.GetNetworkConfigHandler {
	return &getNetworkConfigHandler{configManager: configManager}
}

// handles the request for getting the network configuration.
func (h *getNetworkConfigHandler) Handle(_ operations.GetNetworkConfigParams) middleware.Responder {
	// Build the response with the current network information.
	currentNetwork := h.configManager.CurrentNetwork()
	response := &models.NetworkManagerItem{
		CurrentNetwork:    &currentNetwork.Name,
		AvailableNetworks: *h.configManager.Networks(),
	}

	return operations.NewGetNetworkConfigOK().WithPayload(response)
}
