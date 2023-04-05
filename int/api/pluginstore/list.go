package pluginstore

import (
	"fmt"
	"log"
	"runtime"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/models"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/store"
)

func newList() operations.GetPluginStoreHandler {
	return &list{}
}

type list struct{}

func getDLChecksumAndOs(plugin store.Plugin) (string, string, string, error) {
	pluginURL := ""
	os := runtime.GOOS
	checksum := ""
	switch os { //nolint:varnamelen
	case "linux":
		pluginURL = plugin.Assets.Linux.URL
		checksum = plugin.Assets.Linux.Checksum
	case "darwin":
		switch arch := runtime.GOARCH; arch {
		case "amd64":
			pluginURL = plugin.Assets.MacosAmd64.URL
			checksum = plugin.Assets.MacosAmd64.Checksum
		case "arm64":
			pluginURL = plugin.Assets.MacosArm64.URL
			checksum = plugin.Assets.MacosArm64.Checksum
		default:
			return pluginURL, os, checksum, fmt.Errorf("unsupported OS '%s' and arch '%s'", os, arch)
		}
	case "windows":
		pluginURL = plugin.Assets.Windows.URL
		checksum = plugin.Assets.Windows.Checksum
	default:
		return pluginURL, os, checksum, fmt.Errorf("unsupported OS '%s'", os)
	}

	return pluginURL, os, checksum, nil
}

func (l *list) Handle(_ operations.GetPluginStoreParams) middleware.Responder {
	log.Println("[GET /plugin-store]")

	plugins, err := store.FetchPluginList()
	if err != nil {
		return operations.NewPluginManagerListInternalServerError().WithPayload(
			&models.Error{Code: errorCodeFetchStore, Message: fmt.Sprintf("fetch store plugin list: %s", err.Error())})
	}

	payload := make([]*models.PluginStoreItem, len(plugins))

	for i, plugin := range plugins {
		pluginURL, os, checksum, _ := getDLChecksumAndOs(plugin)
		plugin := plugin
		payload[i] = &models.PluginStoreItem{
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
