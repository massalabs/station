package api

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/onchain/website"
)

func BrowseHandler(params operations.BrowseParams) middleware.Responder {
	body, dynamic, status, err := website.Fetch(node.NewDefaultClient(), params.Address, params.Resource)
	if err != nil {
		if err.Error() == "no data in candidate value key" {
			return NewNotFoundResponder()
		}

		return NewInternalServerErrorResponder(err)
	}

	var contentTypeHeader map[string]string

	if dynamic {
		contentTypeHeader = map[string]string{"Content-Type": "text/html"}
	} else {
		contentTypeHeader = contentType(params.Resource)
	}

	return NewCustomResponder(body, contentTypeHeader, status)
}
