package massa

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/models"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/config"
)

func NewGetNodeHandler(config *config.AppConfig) operations.GetNodeHandler {
	return &getNodeHandler{config: config}
}

type getNodeHandler struct {
	config *config.AppConfig
}

func (h *getNodeHandler) Handle(params operations.GetNodeParams) middleware.Responder {
	return operations.NewGetNodeOK().
		WithPayload(&models.MassaNodeItem{
			Network: h.config.Network,
			URL:     &h.config.NodeURL,
			DNS:     h.config.DNSAddress,
		})
}
