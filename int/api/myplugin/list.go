package myplugin

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/plugin"
)

func list(param operations.PluginManagerListParams) middleware.Responder {
	if manager == nil {
		manager = plugin.NewManager()
	}

	ids := manager.ID()

	payload := make([]*operations.PluginManagerListOKBodyItems0, len(ids))

	for index, id := range ids {
		payload[index] = &operations.PluginManagerListOKBodyItems0{
			ID: ids[index],
		}

		info := manager.Plugin(id).Information()

		if info != nil {
			payload[index].Name = info.Name
			payload[index].Description = info.Description
			payload[index].Logo = info.Logo
		}
	}

	return operations.NewPluginManagerListOK().WithPayload(payload)
}
