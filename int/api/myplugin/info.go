package myplugin

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station/api/swagger/server/models"
	"github.com/massalabs/station/api/swagger/server/restapi/operations"
	"github.com/massalabs/station/pkg/logger"
	"github.com/massalabs/station/pkg/plugin"
)

func newInfo(manager *plugin.Manager) operations.PluginManagerGetInformationHandler {
	return &info{manager: manager}
}

type info struct {
	manager *plugin.Manager
}

func (i *info) Handle(param operations.PluginManagerGetInformationParams) middleware.Responder {
	logger.Debugf("[GET /plugin-manager/%s]", param.ID)

	plgn, err := i.manager.Plugin(param.ID)
	if err != nil {
		return operations.NewPluginManagerGetInformationNotFound().WithPayload(
			&models.Error{Code: errorCodePluginUnknown, Message: "get plugin error: " + err.Error()})
	}

	info := plgn.Information()

	var payload *models.Plugin

	if info != nil {
		pluginURL := fmt.Sprintf("%s%s/", plugin.EndpointPattern, plugin.Alias(info.Author, info.Name))

		payload = &models.Plugin{
			ID:          param.ID,
			Name:        info.Name,
			Author:      info.Author,
			Description: info.Description,
			Logo:        info.Logo,
			Home:        pluginURL + info.Home,
			Updatable:   info.Updatable,
			Version:     info.Version,
			Status:      plgn.Status().String(),
		}
	}

	return operations.NewPluginManagerGetInformationOK().WithPayload(payload)
}
