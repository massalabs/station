package api

import (
	"embed"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/int/api/websites"
)

const indexHTML = "index.html"

const basePath = "html/front/"

const pathFromRoot = "int/api/" + basePath

//go:embed html/front
var content embed.FS

//nolint:nolintlint,ireturn
func ThyraWalletHandler(params operations.ThyraWalletParams) middleware.Responder {
	file := params.Resource
	if params.Resource == indexHTML {
		file = "wallet.html"
	}

	resource, err := content.ReadFile(basePath + file)
	if err != nil {
		return operations.NewThyraWalletNotFound()
	}

	return NewCustomResponder(resource, contentType(params.Resource), http.StatusOK)
}

type WebSiteCreatorData struct {
	UploadMaxSize int
}

//nolint:nolintlint,ireturn
func ThyraWebsiteCreatorHandler(params operations.ThyraWebsiteCreatorParams) middleware.Responder {
	file := params.Resource
	if params.Resource == indexHTML {
		file = "website.html"
		filename := pathFromRoot + file
		maxArchiveSize := websites.GetMaxArchiveSize()

		return NewTemplateResponder(filename, contentType(params.Resource), WebSiteCreatorData{maxArchiveSize})
	}

	resource, err := content.ReadFile(basePath + file)
	if err != nil {
		return operations.NewThyraWebsiteCreatorNotFound()
	}

	return NewCustomResponder(resource, contentType(params.Resource), http.StatusOK)
}

//nolint:nolintlint,ireturn
func ThyraRegistryHandler(params operations.ThyraRegistryParams) middleware.Responder {
	file := params.Resource
	if params.Resource == indexHTML {
		file = "registry.html"
	}

	resource, err := content.ReadFile(basePath + file)
	if err != nil {
		return operations.NewThyraWebsiteCreatorNotFound()
	}

	return NewCustomResponder(resource, contentType(params.Resource), http.StatusOK)
}
