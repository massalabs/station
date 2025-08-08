package api

import (
	"time"

	"github.com/go-openapi/loads"
	"github.com/massalabs/station/api/swagger/server/restapi"
	"github.com/massalabs/station/api/swagger/server/restapi/operations"
	"github.com/massalabs/station/int/api/cmd"
	"github.com/massalabs/station/int/api/massa"
	"github.com/massalabs/station/int/api/myplugin"
	"github.com/massalabs/station/int/api/network"
	"github.com/massalabs/station/int/api/pluginstore"
	"github.com/massalabs/station/int/api/version"
	"github.com/massalabs/station/int/config"
	"github.com/massalabs/station/pkg/logger"
	"github.com/massalabs/station/pkg/plugin"
	"github.com/massalabs/station/pkg/store"
)

type StartServerFlags struct {
	Port    int
	TLSPort int
}

func setAPIFlags(server *restapi.Server, startFlags StartServerFlags) {
	server.Port = startFlags.Port
	server.TLSPort = startFlags.TLSPort
}

func initLocalAPI(
	localAPI *operations.MassastationAPI,
	pluginManager *plugin.Manager,
	configManager *config.MSConfigManager,
) {
	localAPI.CmdExecuteFunctionHandler = cmd.NewExecuteFunctionHandler(configManager)
	localAPI.CmdReadOnlyCallSCHandler = cmd.NewReadOnlyCallSCHandler(configManager)
	localAPI.CmdExecuteSCHandler = cmd.NewExecuteSCHandler(configManager)

	localAPI.MassaGetAddressesHandler = massa.NewGetAddressHandler(configManager)
	localAPI.GetNodeHandler = massa.NewGetNodeHandler(configManager)

	localAPI.GetMassaStationVersionHandler = operations.GetMassaStationVersionHandlerFunc(version.Handle)

	localAPI.CmdDeploySCHandler = cmd.NewDeploySCHandler(configManager)
	localAPI.CmdReadOnlyExecuteSCHandler = cmd.NewReadOnlyExecuteSCHandler(configManager)

	localAPI.EventsGetterHandler = NewEventListenerHandler(configManager)
	localAPI.MassaStationWebAppHandler = operations.MassaStationWebAppHandlerFunc(MassaStationWebAppHandler)

	localAPI.SwitchNetworkHandler = network.NewSwitchNetworkHandler(configManager)
	localAPI.GetNetworkConfigHandler = network.NewGetNetworkConfigHandler(configManager)
	localAPI.CreateNetworkHandler = network.NewCreateNetworkHandler(configManager)
	localAPI.UpdateNetworkHandler = network.NewUpdateNetworkHandler(configManager)
	localAPI.DeleteNetworkHandler = network.NewDeleteNetworkHandler(configManager)

	pluginstore.InitializePluginStoreAPI(localAPI)
	myplugin.InitializePluginAPI(localAPI, pluginManager)
}

type Server struct {
	api           *restapi.Server
	localAPI      *operations.MassastationAPI
	shutdown      chan struct{}
	configManager *config.MSConfigManager
}

// Creates a new server instance and configures it with the given flags.
func NewServer(flags StartServerFlags, configManager *config.MSConfigManager) *Server {
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		logger.Fatal(err.Error())
	}

	localAPI := operations.NewMassastationAPI(swaggerSpec)
	server := restapi.NewServer(localAPI)

	setAPIFlags(server, flags)

	timeout := 5
	server.GracefulTimeout = time.Duration(timeout) * time.Second

	err = store.NewStore()
	if err != nil {
		logger.Error(err.Error())
	}

	return &Server{
		api:           server,
		localAPI:      localAPI,
		shutdown:      make(chan struct{}),
		configManager: configManager,
	}
}

// Starts the server.
// This function starts the server in a new goroutine to avoid blocking the main thread.
func (server *Server) Start(pluginManager *plugin.Manager) {
	initLocalAPI(server.localAPI, pluginManager, server.configManager)
	server.api.ConfigureMassaStationAPI(server.shutdown)

	go func() {
		if err := server.api.Serve(); err != nil {
			logger.Fatal(err.Error())
		}
	}()

	logger.Debug("Server started")
}

// Stops the server and waits for it to finish.
func (server *Server) Stop() {
	logger.Info("Stopping server...")

	if err := server.api.Shutdown(); err != nil {
		logger.Fatal(err.Error())
	}

	<-server.shutdown

	logger.Debug("Server stopped")
}
