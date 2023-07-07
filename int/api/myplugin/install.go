package myplugin

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station/api/swagger/server/models"
	"github.com/massalabs/station/api/swagger/server/restapi/operations"
	"github.com/massalabs/station/pkg/logger"
	"github.com/massalabs/station/pkg/plugin"
)

func newInstall(manager *plugin.Manager) operations.PluginManagerInstallHandler {
	return &install{manager: manager}
}

type install struct {
	manager *plugin.Manager
}

func (i *install) Handle(param operations.PluginManagerInstallParams) middleware.Responder {
	logger.Debugf("[POST /plugin-manager] source: %s", param.Source)

	_, err := url.ParseRequestURI(strings.TrimSpace(param.Source))
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
