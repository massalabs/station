package myplugin

import (
	"log"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/models"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/plugin"
)

func newUpdate(manager *plugin.Manager) operations.PluginManagerUpdateHandler {
	return &update{manager: manager}
}

type update struct {
	manager *plugin.Manager
}

func (u *update) Handle(param operations.PluginManagerUpdateParams) middleware.Responder {
	log.Printf("[PUT /plugin-manager/%s]", param.ID)

	err := u.manager.Update(param.ID)
	if err != nil {
		return operations.NewPluginManagerUpdateInternalServerError().WithPayload(
			&models.Error{Code: "", Message: err.Error()},
		)
	}

	return operations.NewPluginManagerUpdateNoContent()
}
