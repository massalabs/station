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

func (l *list) Handle(param operations.GetPluginStoreParams) middleware.Responder {
	log.Println("[GET /plugin-store]")

	plugins, err := store.FetchPluginList()
	if err != nil {
		return operations.NewPluginManagerListInternalServerError().WithPayload(
			&models.Error{Code: errorCodeFetchStore, Message: fmt.Sprintf("fetch store plugin list: %s", err.Error())})
	}

	payload := make([]*models.PluginStoreItem, len(plugins))

	for i, plugin := range plugins {
		plugin := plugin
		payload[i] = &models.PluginStoreItem{
			Name:        &plugin.Name,
			Description: &plugin.Description,
			Version:     &plugin.Version,
			URL:         &plugin.URL,
			Assets: &models.PluginStoreItemAssets{
				Linux: &models.File{
					URL:      &plugin.Assets.Linux.URL,
					Checksum: &plugin.Assets.Linux.Checksum,
				},
				Windows: &models.File{
					URL:      &plugin.Assets.Windows.URL,
					Checksum: &plugin.Assets.Windows.Checksum,
				},
				MacosAmd64: &models.File{
					URL:      &plugin.Assets.MacosAmd64.URL,
					Checksum: &plugin.Assets.MacosAmd64.Checksum,
				},
				MacosArm64: &models.File{
					URL:      &plugin.Assets.MacosArm64.URL,
					Checksum: &plugin.Assets.MacosArm64.Checksum,
				},
			},
		}
	}

	return operations.NewGetPluginStoreOK().WithPayload(payload)
}
