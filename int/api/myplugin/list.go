package myplugin

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/plugin"
)

//nolint:ireturn
func newList(manager *plugin.Manager) operations.PluginManagerListHandler {
	return &list{manager: manager}
}

type list struct {
	manager *plugin.Manager
}

func (l *list) Handle(param operations.PluginManagerListParams) middleware.Responder {
	ids := l.manager.ID()

	payload := make([]*operations.PluginManagerListOKBodyItems0, len(ids))

	//nolint:varnamelen
	for index, id := range ids {
		//nolint:exhaustruct
		payload[index] = &operations.PluginManagerListOKBodyItems0{
			ID: ids[index],
		}

		info := l.manager.Plugin(id).Information()

		if info != nil {
			payload[index].Name = info.Name
			payload[index].Description = info.Description
			payload[index].Logo = info.Logo
		}
	}

	return operations.NewPluginManagerListOK().WithPayload(payload)
}
