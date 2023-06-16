package myplugin

import (
	"fmt"
	"log"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station/api/swagger/server/models"
	"github.com/massalabs/station/api/swagger/server/restapi/operations"
	"github.com/massalabs/station/pkg/plugin"
)

func newInfo(manager *plugin.Manager) operations.PluginManagerGetInformationHandler {
	return &info{manager: manager}
}

type info struct {
	manager *plugin.Manager
}

func (i *info) Handle(param operations.PluginManagerGetInformationParams) middleware.Responder {
	log.Printf("[GET /plugin-manager/%s]", param.ID)

	plugin, err := i.manager.Plugin(param.ID)
	if err != nil {
		return operations.NewPluginManagerGetInformationNotFound().WithPayload(
			&models.Error{Code: errorCodePluginUnknown, Message: fmt.Sprintf("get plugin error: %s", err.Error())})
	}

	return operations.NewPluginManagerGetInformationOK().WithPayload(
		&operations.PluginManagerGetInformationOKBody{Status: plugin.Status().String()},
	)
}
