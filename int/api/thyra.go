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

const basePathReact = "dist/"

//go:embed html/front
var content embed.FS

//nolint:typecheck,nolintlint
//go:embed dist
var contentReact embed.FS

type WebSiteCreatorData struct {
	UploadMaxSize int
}

func ThyraWebsiteCreatorHandler(params operations.ThyraWebsiteCreatorParams) middleware.Responder {
	file := params.Resource

	if params.Resource == indexHTML {
		file = "website.html"
	}

	resource, err := content.ReadFile(basePath + file)
	if err != nil {
		return operations.NewThyraWebsiteCreatorNotFound()
	}

	if params.Resource == indexHTML {
		maxArchiveSize := websites.GetMaxArchiveSize()

		return NewTemplateResponder(string(resource), contentType(params.Resource), WebSiteCreatorData{maxArchiveSize})
	}

	return NewCustomResponder(resource, contentType(params.Resource), http.StatusOK)
}

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

//nolint:nolintlint,ireturn
func ThyraHomeHandler(params operations.ThyraHomeParams) middleware.Responder {
	content, err := contentReact.ReadFile(basePathReact + "home/" + params.Resource)
	if err != nil {
		return operations.NewThyraHomeNotFound()
	}

	return NewCustomResponder(content, contentType(params.Resource), http.StatusOK)
}

func ThyraPluginManagerHandler(params operations.ThyraPluginManagerParams) middleware.Responder {
	content, err := contentReact.ReadFile(basePathReact + "plugin-manager/" + params.Resource)
	if err != nil {
		return operations.NewThyraPluginManagerNotFound()
	}

	return NewCustomResponder(content, contentType(params.Resource), http.StatusOK)
}
