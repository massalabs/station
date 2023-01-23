package myplugin

import (
	"fmt"
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
	wantedPlugin := r.manager.Plugin(param.Body.ID)
	if wantedPlugin == nil {
		return operations.NewPluginManagerRegisterNotFound()
	}

	urlPlugin, err := url.Parse(param.Body.URL)
	if err != nil {
		return operations.NewPluginManagerRegisterBadRequest().WithPayload(
			&models.Error{Code: errorCodePluginRegisterInvalidData, Message: fmt.Sprintf("parsing URL: %s", err.Error())},
		)
	}

	// Set plugin information.
	info := plugin.Information{
		Name: param.Body.Name, Author: param.Body.Author,
		Description: param.Body.Description,
		Logo:        param.Body.Logo,
		URL:         urlPlugin, APISpec: param.Body.APISpec,
	}

	wantedPlugin.SetInformation(&info)

	// Add alias for http requests.
	alias := fmt.Sprintf("%s/%s", param.Body.Author, param.Body.Name)

	err = r.manager.SetAlias(alias, param.Body.ID)
	if err != nil {
		return operations.NewPluginManagerRegisterBadRequest().WithPayload(
			&models.Error{Code: errorCodePluginRegisterUnknown, Message: fmt.Sprintf("setting alias: %s", err.Error())},
		)
	}

	return operations.NewPluginManagerRegisterNoContent()
}
