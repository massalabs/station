package api

import (
	"embed"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station/api/swagger/server/restapi/operations"
	"github.com/massalabs/station/int/api/utils"
)

const basePathReact = "dist/"

//nolint:typecheck,nolintlint
//go:embed dist
var contentReact embed.FS

type WebSiteCreatorData struct {
	UploadMaxSize int
}

func MassaStationWebAppHandler(params operations.MassaStationWebAppParams) middleware.Responder {
	resourceName := params.Resource

	resourceContent, err := contentReact.ReadFile(basePathReact + "massastation/" + resourceName)
	if err != nil {
		resourceName = "index.html"
		resourceContent, err = contentReact.ReadFile(basePathReact + "massastation/" + resourceName)
		if err != nil {
			return operations.NewMassaStationWebAppNotFound()
		}
	}

	return utils.NewCustomResponder(resourceContent, utils.ContentType(resourceName), http.StatusOK)
}
