package pluginstore

import (
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/store"
)

func InitializePluginStoreAPI(api *operations.ThyraServerAPI) {
	api.GetPluginStoreHandler = newList(store.StoreInstance)
}

const (
	errorCodeFetchStore = "FetchPluginStore-001"
)
