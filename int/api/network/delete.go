package network

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station/api/swagger/server/models"
	"github.com/massalabs/station/api/swagger/server/restapi/operations"
	"github.com/massalabs/station/int/config"
)

type deleteNetworkHandler struct{ configManager *config.MSConfigManager }

func NewDeleteNetworkHandler(configManager *config.MSConfigManager) operations.DeleteNetworkHandler {
	return &deleteNetworkHandler{configManager: configManager}
}

func (h *deleteNetworkHandler) Handle(params operations.DeleteNetworkParams) middleware.Responder {
	if err := h.configManager.DeleteNetwork(params.Network); err != nil {
		if err.Error() == "unknown network: "+params.Network {
			return operations.NewDeleteNetworkNotFound().WithPayload(&models.Error{Code: errorCodeNetworkUnknown, Message: err.Error()})
		}
		return operations.NewDeleteNetworkBadRequest().WithPayload(&models.Error{Code: errorCodeNetworkDeleteFailed, Message: err.Error()})
	}

	current := h.configManager.CurrentNetwork()
	response := &models.NetworkManagerItem{CurrentNetwork: &current.Name, AvailableNetworks: *h.configManager.Networks()}
	return operations.NewDeleteNetworkOK().WithPayload(response)
}
