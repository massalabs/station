package api

import (
	"embed"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
)

//go:embed html/front
var content embed.FS

//nolint:nolintlint,ireturn
func ThyraWalletHandler(params operations.ThyraWalletParams) middleware.Responder {
	basePath := "html/front/"

	file := params.Resource
	if params.Resource == "index.html" {
		file = "wallet.html"
	}

	resource, err := content.ReadFile(basePath + file)
	if err != nil {
		return operations.NewThyraWalletNotFound()
	}

	return NewCustomResponder(resource, contentType(params.Resource), http.StatusOK)
}

//nolint:nolintlint,ireturn
func ThyraWebsiteCreatorHandler(params operations.ThyraWebsiteCreatorParams) middleware.Responder {
	basePath := "html/front/"

	file := params.Resource
	if params.Resource == "index.html" {
		file = "website.html"
	}

	resource, err := content.ReadFile(basePath + file)
	if err != nil {
		return operations.NewThyraWebsiteCreatorNotFound()
	}

	return NewCustomResponder(resource, contentType(params.Resource), http.StatusOK)
}
