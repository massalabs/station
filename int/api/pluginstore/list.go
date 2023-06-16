package pluginstore

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station/api/swagger/server/models"
	"github.com/massalabs/station/api/swagger/server/restapi/operations"
	"github.com/massalabs/station/pkg/store"
)

func newList() operations.GetPluginStoreHandler {
	return &list{store: nil}
}

type list struct {
	store *store.Store
}

func (l *list) Handle(_ operations.GetPluginStoreParams) middleware.Responder {
	plugins := store.StoreInstance.Plugins

	payload := make([]*models.PluginStoreItem, len(plugins))

	for index, plugin := range plugins {
		//nolint:varnamelen
		pluginURL, os, checksum, err := plugin.GetDLChecksumAndOs()
		if err != nil {
			return operations.NewPluginManagerListInternalServerError().WithPayload(
				&models.Error{Code: errorCodeFetchStore, Message: fmt.Sprintf("Error getting OS info: %s", err.Error())})
		}

		plugin := plugin
		payload[index] = &models.PluginStoreItem{
			Name:        &plugin.Name,
			Description: &plugin.Description,
			Version:     &plugin.Version,
			URL:         &plugin.URL,
			File: &models.File{
				URL:      &pluginURL,
				Checksum: &checksum,
			},
			Os: os,
		}
	}

	return operations.NewGetPluginStoreOK().WithPayload(payload)
}
