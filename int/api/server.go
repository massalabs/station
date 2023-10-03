package api

import (
	"fmt"
	"time"

	"github.com/go-openapi/loads"
	"github.com/massalabs/station/api/swagger/server/restapi"
	"github.com/massalabs/station/api/swagger/server/restapi/operations"
	"github.com/massalabs/station/int/api/cmd"
	"github.com/massalabs/station/int/api/massa"
	"github.com/massalabs/station/int/api/myplugin"
	"github.com/massalabs/station/int/api/network"
	"github.com/massalabs/station/int/api/pluginstore"
	"github.com/massalabs/station/int/api/websites"
	"github.com/massalabs/station/int/config"
	"github.com/massalabs/station/pkg/logger"
	"github.com/massalabs/station/pkg/node"
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
	networkManager *config.NetworkManager,
	pluginManager *plugin.Manager,
) {
	config := networkManager.Network()

	localAPI.CmdExecuteFunctionHandler = cmd.NewExecuteFunctionHandler(config)

	localAPI.MassaGetAddressesHandler = massa.NewGetAddressHandler(config)
	localAPI.GetNodeHandler = massa.NewGetNodeHandler(config)

	localAPI.GetMassaStationVersionHandler = operations.GetMassaStationVersionHandlerFunc(massa.GetMassaStationVersionFunc)

	localAPI.CmdDeploySCHandler = cmd.NewDeploySCHandler(config)

	localAPI.WebsiteUploaderPrepareHandler = websites.NewWebsitePrepareHandler(config)
	localAPI.WebsiteUploaderUploadHandler = websites.NewWebsiteUploadHandler(config)
	localAPI.WebsiteUploadMissingChunksHandler = websites.NewWebsiteUploadMissedChunkHandler(config)

	localAPI.MyDomainsGetterHandler = websites.NewDomainsHandler(config)
	localAPI.AllDomainsGetterHandler = websites.NewRegistryHandler(config)

	localAPI.EventsGetterHandler = NewEventListenerHandler(config)
	localAPI.MassaStationWebAppHandler = operations.MassaStationWebAppHandlerFunc(MassaStationWebAppHandler)

	localAPI.SwitchNetworkHandler = network.NewSwitchNetworkHandler(networkManager)
	localAPI.GetNetworkConfigHandler = network.NewGetNetworkConfigHandler(networkManager)

	pluginstore.InitializePluginStoreAPI(localAPI)
	myplugin.InitializePluginAPI(localAPI, pluginManager)
}

type Server struct {
	api      *restapi.Server
	localAPI *operations.MassastationAPI
	shutdown chan struct{}
}

// Creates a new server instance and configures it with the given flags.
func NewServer(flags StartServerFlags) *Server {
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
		api:      server,
		localAPI: localAPI,
		shutdown: make(chan struct{}),
	}
}

// Starts the server.
// This function starts the server in a new goroutine to avoid blocking the main thread.
func (server *Server) Start(networkManager *config.NetworkManager, pluginManager *plugin.Manager) {
	server.printNodeVersion(networkManager)

	initLocalAPI(server.localAPI, networkManager, pluginManager)
	server.api.ConfigureMassaStationAPI(*networkManager.Network(), server.shutdown)

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

// Displays the node version of the connected node.
func (server *Server) printNodeVersion(networkManager *config.NetworkManager) {
	client := node.NewClient(networkManager.Network().NodeURL)
	status, err := node.Status(client)

	if err == nil {
		logger.Info(
			fmt.Sprintf("Connected to node server %s (version %s)",
				networkManager.Network().NodeURL, *status.Version),
		)
	} else {
		logger.Errorf("Could not get node version: %s", err.Error())
	}
}
