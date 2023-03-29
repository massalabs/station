package pluginstore

import (
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
)

func InitializePluginStoreAPI(api *operations.ThyraServerAPI) {
	api.GetPluginStoreHandler = newList()
}

const (
	errorCodeFetchStore = "FetchPluginStore-001"
)
