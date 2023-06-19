package pluginstore

import (
	"github.com/massalabs/station/api/swagger/server/restapi/operations"
)

func InitializePluginStoreAPI(api *operations.MassastationAPI) {
	api.GetPluginStoreHandler = newList()
}

const (
	errorCodeFetchStore = "FetchPluginStore-001"
)
