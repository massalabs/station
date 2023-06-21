package myplugin

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station/api/swagger/server/models"
	"github.com/massalabs/station/api/swagger/server/restapi/operations"
	"github.com/massalabs/station/pkg/plugin"
)

func newList(manager *plugin.Manager) operations.PluginManagerListHandler {
	return &list{manager: manager}
}

type list struct {
	manager *plugin.Manager
}

func (l *list) Handle(_ operations.PluginManagerListParams) middleware.Responder {
	ids := l.manager.ID()

	payload := make([]*models.Plugin, len(ids))

	for index, id := range ids {

		plgn, err := l.manager.Plugin(id)
		if err != nil {
			return operations.NewPluginManagerListNotFound().WithPayload(
				&models.Error{Code: errorCodePluginUnknown, Message: fmt.Sprintf("get plugin error: %s", err.Error())})
		}

		info := plgn.Information()

		if info != nil {
			pluginURL := fmt.Sprintf("%s%s/", plugin.EndpointPattern, plugin.Alias(info.Author, info.Name))

			payload[index] = &models.Plugin{
				ID:          id,
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
	}
	return operations.NewPluginManagerListOK().WithPayload(payload)
}
