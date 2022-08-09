package api

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/onchain/website"
)

func BrowseHandler(params operations.BrowseParams) middleware.Responder {
	body, err := website.Fetch(node.NewDefaultClient(), params.Address, params.Resource)
	if err != nil {
		if err.Error() == "no data in candidate value key" {
			return NewNotFoundResponder()
		} else {
			return NewInternalServerErrorResponder(err)
		}
	}

	return NewCustomResponder(body, contentType(params.Resource), http.StatusOK)
}
