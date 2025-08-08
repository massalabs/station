package network

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station/api/swagger/server/models"
	"github.com/massalabs/station/api/swagger/server/restapi/operations"
	"github.com/massalabs/station/int/config"
)

type createNetworkHandler struct{ configManager *config.MSConfigManager }

func NewCreateNetworkHandler(configManager *config.MSConfigManager) operations.CreateNetworkHandler {
	return &createNetworkHandler{configManager: configManager}
}

func (h *createNetworkHandler) Handle(params operations.CreateNetworkParams) middleware.Responder {
	body := params.Body

	// Required fields are enforced by swagger validation; just read values
	makeDefault := false
	if body.Default != nil {
		makeDefault = *body.Default
	}

	if err := h.configManager.AddNetwork(body.Name, body.URL, makeDefault); err != nil {
		return operations.NewCreateNetworkBadRequest().WithPayload(&models.Error{Code: errorCodeNetworkCreateFailed, Message: err.Error()})
	}

	current := h.configManager.CurrentNetwork()
	response := &models.NetworkManagerItem{CurrentNetwork: &current.Name, AvailableNetworkInfos: []*models.NetworkInfoItem{}}
	return operations.NewCreateNetworkOK().WithPayload(response)
}
