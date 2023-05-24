package pluginstore

import (
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/config"
)

func InitializePluginStoreAPI(api *operations.ThyraServerAPI, config *config.AppConfig) {
	api.GetPluginStoreHandler = newList(config.Store)
}

const (
	errorCodeFetchStore = "FetchPluginStore-001"
)
