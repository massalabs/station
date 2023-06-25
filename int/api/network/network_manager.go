package network

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station/api/swagger/server/models"
	"github.com/massalabs/station/api/swagger/server/restapi/operations"
	"github.com/massalabs/station/pkg/config"
)

func NewSwitchNetworkHandler(networkManager *config.NetworkManager) operations.SwitchNetworkHandler {
	return &switchNetworkHandler{networkManager: networkManager}
}

type switchNetworkHandler struct {
	networkManager *config.NetworkManager
}

func (h *switchNetworkHandler) Handle(params operations.SwitchNetworkParams) middleware.Responder {
	err := h.networkManager.SwitchNetwork(params.Network)
	if err != nil {
		return operations.NewSwitchNetworkNotFound().WithPayload(
			&models.Error{
				Code:    "404",
				Message: "Network not found",
			},
		)
	}

	response := &models.NetworkManagerItem{
		ActualNetwork:     &h.networkManager.Network().Network,
		AvailableNetworks: *h.networkManager.Networks(),
	}
	return operations.NewSwitchNetworkOK().WithPayload(response)
}

func NewGetNetworkConfigHandler(networkManager *config.NetworkManager) operations.GetNetworkConfigHandler {
	return &getNetworkConfigHandler{networkManager: networkManager}
}

type getNetworkConfigHandler struct {
	networkManager *config.NetworkManager
}

func (h *getNetworkConfigHandler) Handle(params operations.GetNetworkConfigParams) middleware.Responder {
	response := &models.NetworkManagerItem{
		ActualNetwork:     &h.networkManager.Network().Network,
		AvailableNetworks: *h.networkManager.Networks(),
	}

	return operations.NewGetNetworkConfigOK().WithPayload(response)
}
