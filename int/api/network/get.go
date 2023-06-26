package network

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station/api/swagger/server/models"
	"github.com/massalabs/station/api/swagger/server/restapi/operations"
	"github.com/massalabs/station/pkg/config"
)

type getNetworkConfigHandler struct {
	networkManager *config.NetworkManager
}

// NewGetNetworkConfigHandler creates a new getNetworkConfigHandler instance.
func NewGetNetworkConfigHandler(networkManager *config.NetworkManager) operations.GetNetworkConfigHandler {
	return &getNetworkConfigHandler{networkManager: networkManager}
}

// handles the request for getting the network configuration.
func (h *getNetworkConfigHandler) Handle(_ operations.GetNetworkConfigParams) middleware.Responder {
	// Build the response with the current network information.
	response := &models.NetworkManagerItem{
		ActualNetwork:     &h.networkManager.Network().Network,
		AvailableNetworks: *h.networkManager.Networks(),
	}

	return operations.NewGetNetworkConfigOK().WithPayload(response)
}
