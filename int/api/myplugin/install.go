package myplugin

import (
	"fmt"
	"net/url"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/models"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/plugin"
)

func newInstall(manager *plugin.Manager) operations.PluginManagerInstallHandler {
	return &install{manager: manager}
}

type install struct {
	manager *plugin.Manager
}

func (i *install) Handle(param operations.PluginManagerInstallParams) middleware.Responder {
	_, err := url.ParseRequestURI(param.Source)
	if err != nil {
		return operations.NewPluginManagerInstallBadRequest().WithPayload(
			&models.Error{
				Code:    errorCodePluginInstallationInvalidSource,
				Message: fmt.Sprintf("Error: given source %s is not a valid URL (%s)", param.Source, err),
			})
	}

	err = i.manager.Install(param.Source)
	if err != nil {
		return operations.NewPluginManagerInstallInternalServerError().WithPayload(
			&models.Error{
				Code:    errorCodePluginUnknown,
				Message: fmt.Sprintf("Error: %s", err),
			})
	}

	return operations.NewPluginManagerInstallNoContent()
}
