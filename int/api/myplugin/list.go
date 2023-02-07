package myplugin

import (
	"fmt"
	"log"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/models"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/plugin"
)

func newList(manager *plugin.Manager) operations.PluginManagerListHandler {
	return &list{manager: manager}
}

type list struct {
	manager *plugin.Manager
}

func (l *list) Handle(param operations.PluginManagerListParams) middleware.Responder {
	log.Println("[GET /plugin-manager]")

	ids := l.manager.ID()

	payload := make([]*operations.PluginManagerListOKBodyItems0, len(ids))

	//nolint:varnamelen
	for index, id := range ids {
		//nolint:exhaustruct
		payload[index] = &operations.PluginManagerListOKBodyItems0{
			ID: ids[index],
		}

		plgn, err := l.manager.Plugin(id)
		if err != nil {
			return operations.NewPluginManagerListNotFound().WithPayload(
				&models.Error{Code: errorCodePluginUnknown, Message: fmt.Sprintf("get plugin error: %s", err.Error())})
		}

		info := plgn.Information()

		if info != nil {
			payload[index].Name = info.Name
			payload[index].Description = info.Description
			payload[index].Logo = info.Logo
		}
	}

	return operations.NewPluginManagerListOK().WithPayload(payload)
}
