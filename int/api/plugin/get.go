package plugin

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/models"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	pluginManager "github.com/massalabs/thyra/pkg/plugins"
)

//nolint:nolintlint,ireturn
func NewGet(manager *pluginManager.PluginManager) operations.MgmtPluginsListHandler {
	return &pluginHandler{manager: manager}
}

type pluginHandler struct {
	manager *pluginManager.PluginManager
}

//nolint:nolintlint,ireturn
func (pluginCatalog *pluginHandler) Handle(params operations.MgmtPluginsListParams) middleware.Responder {
	pluginList := pluginCatalog.manager.List()

	var body []*models.Plugin

	for i := 0; i < len(pluginList); i++ {
		pluginInfo := &models.Plugin{
			Name: pluginList[i].Manifest.Name,
			Port: int64(pluginList[i].Port),
		}

		body = append(body, pluginInfo)
	}

	return operations.NewMgmtPluginsListOK().WithPayload(body)
}
