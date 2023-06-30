package myplugin

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station/api/swagger/server/models"
	"github.com/massalabs/station/api/swagger/server/restapi/operations"
	"github.com/massalabs/station/int/api/utils"
	"github.com/massalabs/station/pkg/config"
	"github.com/massalabs/station/pkg/plugin"
)

func newLogo(manager *plugin.Manager) operations.PluginManagerGetLogoHandler {
	return &logo{manager: manager}
}

type logo struct {
	manager *plugin.Manager
}

func (l *logo) Handle(param operations.PluginManagerGetLogoParams) middleware.Responder {
	plgn, err := l.manager.Plugin(param.ID)
	if err != nil {
		return operations.NewPluginManagerGetLogoNotFound().WithPayload(
			&models.Error{Code: errorCodePluginUnknown, Message: fmt.Sprintf("plugin %s not found", param.ID)})
	}

	info := plgn.Information()
	logoPath := info.Logo

	// Open the logo file
	logoFile, err := os.Open(logoPath)
	if err != nil {
		config.Logger.Errorf("Error opening logo file: %s", err)

		if err != nil {
			return operations.NewPluginManagerExecuteCommandNotFound().WithPayload(
				&models.Error{Code: errorCodePluginLogoNotFound, Message: fmt.Sprintf("plugin %s logo not found", param.ID)})
		}
	}
	defer logoFile.Close()

	// Read the logo file
	logoData, err := ioutil.ReadAll(logoFile)
	if err != nil {
		config.Logger.Errorf("Error reading logo file: %s", err)

		return operations.NewPluginManagerGetLogoInternalServerError().WithPayload(
			&models.Error{Code: errorCodePluginUnknown, Message: fmt.Sprintf("get plugin logo error: %s", err.Error())})
	}

	return utils.NewCustomResponder(logoData, utils.ContentType(logoPath), http.StatusOK)
}
