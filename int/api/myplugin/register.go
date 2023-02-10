package myplugin

import (
	"fmt"
	"log"
	"net/url"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/models"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/plugin"
)

func newRegister(manager *plugin.Manager) operations.PluginManagerRegisterHandler {
	return &register{manager: manager}
}

type register struct {
	manager *plugin.Manager
}

func (r *register) Handle(param operations.PluginManagerRegisterParams) middleware.Responder {
	log.Printf("[POST /plugin-manager/register] Name: %s ID:%s", param.Body.Name, param.Body.ID)

	wantedPlugin, err := r.manager.Plugin(param.Body.ID)
	if err != nil {
		return operations.NewPluginManagerRegisterNotFound().WithPayload(
			&models.Error{Code: errorCodePluginUnknown, Message: fmt.Sprintf("get plugin error: %s", err.Error())})
	}

	urlPlugin, err := url.Parse(param.Body.URL)
	if err != nil {
		return operations.NewPluginManagerRegisterBadRequest().WithPayload(
			&models.Error{Code: errorCodePluginRegisterInvalidData, Message: fmt.Sprintf("parsing Plugin URL: %s", err.Error())},
		)
	}

	// Set plugin information.
	info := plugin.Information{
		Name: param.Body.Name, Author: param.Body.Author,
		Description: param.Body.Description,
		Logo:        param.Body.Logo,
		URL:         urlPlugin,
		APISpec:     param.Body.APISpec,
		Home:        param.Body.Home,
	}

	wantedPlugin.SetInformation(&info)

	// Add alias for http requests.

	alias := fmt.Sprintf("%s/%s", param.Body.Author, param.Body.Name)

	err = r.manager.SetAlias(alias, param.Body.ID)

	if err != nil {
		log.Printf("setting plugin alias: %s", err)

		return operations.NewPluginManagerRegisterBadRequest().WithPayload(
			&models.Error{Code: errorCodePluginRegisterUnknown, Message: fmt.Sprintf("setting alias: %s", err.Error())},
		)
	}

	return operations.NewPluginManagerRegisterNoContent()
}
