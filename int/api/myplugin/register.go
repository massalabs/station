package myplugin

import (
	"fmt"
	"net/url"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station/api/swagger/server/models"
	"github.com/massalabs/station/api/swagger/server/restapi/operations"
	"github.com/massalabs/station/pkg/logger"
	"github.com/massalabs/station/pkg/plugin"
)

func newRegister(manager *plugin.Manager) operations.PluginManagerRegisterHandler {
	return &register{manager: manager}
}

type register struct {
	manager *plugin.Manager
}

func (r *register) Handle(param operations.PluginManagerRegisterParams) middleware.Responder {
	wantedPlugin, err := r.manager.Plugin(param.Body.ID)
	if err != nil {
		return operations.NewPluginManagerRegisterNotFound().WithPayload(
			&models.Error{Code: errorCodePluginUnknown, Message: fmt.Sprintf("get plugin error: %s", err.Error())})
	}

	alias := plugin.Alias(wantedPlugin.Information().Author, wantedPlugin.Information().Name)

	_, err = r.manager.PluginByAlias(alias)
	if err == nil {
		return operations.NewPluginManagerRegisterBadRequest().WithPayload(
			&models.Error{
				Code:    errorCodePluginRegisterAlreadyRegistered,
				Message: fmt.Sprintf("plugin already registered: %s", alias),
			})
	}

	urlPlugin, err := url.Parse(param.Body.URL)
	if err != nil {
		return operations.NewPluginManagerRegisterBadRequest().WithPayload(
			&models.Error{Code: errorCodePluginRegisterInvalidData, Message: fmt.Sprintf("parsing Plugin URL: %s", err.Error())},
		)
	}

	err = wantedPlugin.SetInformation(urlPlugin)
	if err != nil {
		return operations.NewPluginManagerRegisterBadRequest().WithPayload(
			&models.Error{Code: errorCodePluginRegisterInvalidData, Message: fmt.Sprintf("parsing Plugin URL: %s", err.Error())},
		)
	}

	wantedPlugin.InitReverseProxy()

	// Add alias for http requests.

	err = r.manager.SetAlias(alias, param.Body.ID)
	if err != nil {
		logger.Debugf("setting plugin alias: %s", err)

		return operations.NewPluginManagerRegisterBadRequest().WithPayload(
			&models.Error{Code: errorCodePluginRegisterUnknown, Message: fmt.Sprintf("setting alias: %s", err.Error())},
		)
	}

	return operations.NewPluginManagerRegisterNoContent()
}
