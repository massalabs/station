package plugins

import (
	"strconv"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/models"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	pluginManager "github.com/massalabs/thyra/pkg/plugins"
)

//nolint:nolintlint,ireturn
func NewGet(manager *pluginManager.PluginManager) operations.MgmtPluginsListHandler {
	return &pluginsGet{manager: manager}
}

type pluginsGet struct {
	manager *pluginManager.PluginManager
}

//nolint:nolintlint,ireturn
func (c *pluginsGet) Handle(params operations.MgmtPluginsListParams) middleware.Responder {
	pluginList := c.manager.List()

	var plugin []*models.Plugin

	for i := 0; i < len(pluginList); i++ {
		pluginInfo := &models.Plugin{
			Name: pluginList[i].Manifest.Name,
			Port: strconv.Itoa(pluginList[i].Port),
		}

		plugin = append(plugin, pluginInfo)
	}

	return operations.NewMgmtPluginsListOK().WithPayload(plugin)
}
