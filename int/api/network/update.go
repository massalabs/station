package network

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station/api/swagger/server/models"
	"github.com/massalabs/station/api/swagger/server/restapi/operations"
	"github.com/massalabs/station/int/config"
)

type updateNetworkHandler struct{ configManager *config.MSConfigManager }

func NewUpdateNetworkHandler(configManager *config.MSConfigManager) operations.UpdateNetworkHandler {
	return &updateNetworkHandler{configManager: configManager}
}

func (h *updateNetworkHandler) Handle(params operations.UpdateNetworkParams) middleware.Responder {
	body := params.Body

	makeDefault := body.Default

	newURL := body.URL

	newName := body.NewName

	if err := h.configManager.EditNetwork(params.Network, newURL, makeDefault, newName); err != nil {
		if err.Error() == "unknown network: "+params.Network {
			return operations.NewUpdateNetworkNotFound().WithPayload(&models.Error{Code: errorCodeNetworkUnknown, Message: err.Error()})
		}
		return operations.NewUpdateNetworkBadRequest().WithPayload(&models.Error{Code: errorCodeNetworkUpdateFailed, Message: err.Error()})
	}

	current := h.configManager.CurrentNetwork()
	response := &models.NetworkManagerItem{CurrentNetwork: &current.Name, AvailableNetworks: *h.configManager.Networks()}
	return operations.NewUpdateNetworkOK().WithPayload(response)
}
