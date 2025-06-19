package pluginstore

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gosimple/slug"
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
			return operations.NewGetPluginStoreInternalServerError().WithPayload(
				&models.Error{Code: errorCodeFetchStore, Message: "Error getting OS info: " + err.Error()})
		}

		folderLogo := slug.Make(plugin.Name)
		plugin.Logo = fmt.Sprintf(
			"https://massa-station-assets.s3.eu-west-3.amazonaws.com/plugins/icons/%s/%s",
			folderLogo,
			plugin.Logo,
		)
		plugin := plugin
		payload[index] = &models.PluginStoreItem{
			Name:                   &plugin.Name,
			Author:                 plugin.Author,
			Description:            &plugin.Description,
			Version:                &plugin.Version,
			Logo:                   plugin.Logo,
			MassastationMinVersion: plugin.MassaStationVersion,
			URL:                    &plugin.URL,
			File: &models.File{
				URL:      &pluginURL,
				Checksum: &checksum,
			},
			Os:           os,
			IsCompatible: plugin.IsCompatible,
		}
	}

	return operations.NewGetPluginStoreOK().WithPayload(payload)
}
