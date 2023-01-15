package myplugin

import (
	"fmt"
	"net/url"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/models"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/plugin"
)

var manager *plugin.Manager

func pluginInstall(param operations.PluginManagerInstallParams) middleware.Responder {
	if manager == nil {
		manager = plugin.NewManager()
	}

	_, err := url.ParseRequestURI(param.Source)
	if err != nil {
		return operations.NewPluginManagerInstallBadRequest().WithPayload(
			&models.Error{
				Code:    errorCodePluginInstallationInvalidSource,
				Message: fmt.Sprintf("Error: given source %s is not a valid URL (%s)", param.Source, err),
			})
	}

	err = manager.Install(param.Source)
	if err != nil {
		return operations.NewPluginManagerInstallInternalServerError().WithPayload(
			&models.Error{
				Code:    errorCodePluginUnknown,
				Message: fmt.Sprintf("Error: %s", err),
			})
	}

	return operations.NewPluginManagerInstallNoContent()
}
