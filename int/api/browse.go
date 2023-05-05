package api

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/config"
	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/onchain/website"
)

func NewBrowseHandler(config *config.AppConfig) operations.BrowseHandler {
	return &browseHandler{config: config}
}

type browseHandler struct {
	config *config.AppConfig
}

func (h *browseHandler) Handle(params operations.BrowseParams) middleware.Responder {
	body, err := website.Fetch(node.NewClient(h.config.NodeURL), params.Address, params.Resource)
	if err != nil {
		if err.Error() == "no data in candidate value key" {
			return NewNotFoundResponder()
		}

		return NewInternalServerErrorResponder(err)
	}

	return NewCustomResponder(body, contentType(params.Resource), http.StatusOK)
}
