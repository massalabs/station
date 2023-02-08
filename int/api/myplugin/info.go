package myplugin

import (
	"fmt"
	"log"
	"strconv"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/models"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/plugin"
)

func newInfo(manager *plugin.Manager) operations.PluginManagerGetInformationHandler {
	return &info{manager: manager}
}

type info struct {
	manager *plugin.Manager
}

func (i *info) Handle(param operations.PluginManagerGetInformationParams) middleware.Responder {
	log.Printf("[GET /plugin-manager/%s]", param.ID)

	pluginID, err := strconv.ParseInt(param.ID, 10, 64)
	if err != nil {
		return operations.NewPluginManagerGetInformationBadRequest().WithPayload(
			&models.Error{Code: "", Message: err.Error()},
		)
	}

	plugin, err := i.manager.Plugin(pluginID)
	if err != nil {
		return operations.NewPluginManagerGetInformationNotFound().WithPayload(
			&models.Error{Code: errorCodePluginUnknown, Message: fmt.Sprintf("get plugin error: %s", err.Error())})
	}

	return operations.NewPluginManagerGetInformationOK().WithPayload(
		&operations.PluginManagerGetInformationOKBody{Status: plugin.Status().String()},
	)
}
