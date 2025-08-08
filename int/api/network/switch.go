package network

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station/api/swagger/server/models"
	"github.com/massalabs/station/api/swagger/server/restapi/operations"
	"github.com/massalabs/station/int/config"
)

type switchNetworkHandler struct{ configManager *config.MSConfigManager }

func NewSwitchNetworkHandler(configManager *config.MSConfigManager) operations.SwitchNetworkHandler {
	return &switchNetworkHandler{configManager: configManager}
}

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

	currentNetwork := h.configManager.CurrentNetwork()
	availableNetworks := buildAvailableNetworkInfos(h.configManager)
	response := &models.NetworkManagerItem{
		CurrentNetwork:        &currentNetwork.Name,
		AvailableNetworkInfos: availableNetworks,
	}

	return operations.NewSwitchNetworkOK().WithPayload(response)
}
