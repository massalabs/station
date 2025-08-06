package network

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station/api/swagger/server/models"
	"github.com/massalabs/station/api/swagger/server/restapi/operations"
	"github.com/massalabs/station/int/config"
)

const errorCodeNetworkUnknown = "Network-0001"

type switchNetworkHandler struct{ configManager *config.MSConfigManager }

// NewSwitchNetworkHandler creates a new switchNetworkHandler instance.
func NewSwitchNetworkHandler(configManager *config.MSConfigManager) operations.SwitchNetworkHandler {
	return &switchNetworkHandler{configManager: configManager}
}

// handles the request for switching the network.
func (h *switchNetworkHandler) Handle(params operations.SwitchNetworkParams) middleware.Responder {
	err := h.configManager.SwitchNetwork(params.Network)
	if err != nil {
		// If the network is not found, return a 404 response with an error message.
		return operations.NewSwitchNetworkNotFound().WithPayload(
			&models.Error{
				Code:    errorCodeNetworkUnknown,
				Message: "Network not found: " + err.Error(),
			},
		)
	}

	// Build the response with the current network information.
	currentNetwork := h.configManager.CurrentNetwork()
	response := &models.NetworkManagerItem{
		CurrentNetwork:    &currentNetwork.Name,
		AvailableNetworks: *h.configManager.Networks(),
	}

	return operations.NewSwitchNetworkOK().WithPayload(response)
}
