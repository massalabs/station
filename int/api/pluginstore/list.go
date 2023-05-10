package pluginstore

import (
	"fmt"
	"log"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/models"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/store"
)

func newList() operations.GetPluginStoreHandler {
	return &list{}
}

type list struct{}

func (l *list) Handle(_ operations.GetPluginStoreParams) middleware.Responder {
	log.Println("[GET /plugin-store]")

	plugins, err := store.FetchPluginList()
	if err != nil {
		return operations.NewPluginManagerListInternalServerError().WithPayload(
			&models.Error{Code: errorCodeFetchStore, Message: fmt.Sprintf("fetch store plugin list: %s", err.Error())})
	}

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
