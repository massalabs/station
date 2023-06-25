package api

import (
	"log"
	"os"

	"github.com/go-openapi/loads"
	"github.com/jessevdk/go-flags"
	"github.com/massalabs/station/api/swagger/server/restapi"
	"github.com/massalabs/station/api/swagger/server/restapi/operations"
	"github.com/massalabs/station/int/api/cmd"
	"github.com/massalabs/station/int/api/massa"
	"github.com/massalabs/station/int/api/myplugin"
	"github.com/massalabs/station/int/api/pluginstore"
	"github.com/massalabs/station/int/api/websites"
	"github.com/massalabs/station/pkg/config"
	"github.com/massalabs/station/pkg/node"
	"github.com/massalabs/station/pkg/store"
)

type StartServerFlags struct {
	Port              int
	TLSPort           int
	TLSCertificate    string
	TLSCertificateKey string
	MassaNodeServer   string
	DNSAddress        string
	Version           bool
}

func setAPIFlags(server *restapi.Server, startFlags StartServerFlags) {
	server.Port = startFlags.Port
	server.TLSPort = startFlags.TLSPort

	if _, err := os.Stat(startFlags.TLSCertificate); err == nil {
		server.TLSCertificate = flags.Filename(startFlags.TLSCertificate)
	}

	if _, err := os.Stat(startFlags.TLSCertificateKey); err == nil {
		server.TLSCertificateKey = flags.Filename(startFlags.TLSCertificateKey)
	}
}

func initLocalAPI(localAPI *operations.MassastationAPI, networkManager *config.NetworkManager) {
	config := networkManager.Network()

	localAPI.CmdExecuteFunctionHandler = cmd.NewExecuteFunctionHandler(config)

	localAPI.MassaGetAddressesHandler = massa.NewGetAddressHandler(config)
	localAPI.GetNodeHandler = massa.NewGetNodeHandler(config)

	localAPI.CmdDeploySCHandler = cmd.NewDeploySCHandler(config)

	localAPI.WebsiteUploaderPrepareHandler = websites.NewWebsitePrepareHandler(config)
	localAPI.WebsiteUploaderUploadHandler = websites.NewWebsiteUploadHandler(config)
	localAPI.WebsiteUploadMissingChunksHandler = websites.NewWebsiteUploadMissedChunkHandler(config)

	localAPI.MyDomainsGetterHandler = websites.NewDomainsHandler(config)
	localAPI.AllDomainsGetterHandler = websites.NewRegistryHandler(config)

	localAPI.WebOnChainSearchHandler = operations.WebOnChainSearchHandlerFunc(WebOnChainSearchHandler)
	localAPI.MassaStationHomeHandler = operations.MassaStationHomeHandlerFunc(MassaStationHomeHandler)
	localAPI.EventsGetterHandler = NewEventListenerHandler(config)
	localAPI.MassaStationPluginManagerHandler = operations.MassaStationPluginManagerHandlerFunc(MassaStationPluginManagerHandler)
	localAPI.MassaStationWebAppHandler = operations.MassaStationWebAppHandlerFunc(MassaStationWebAppHandler)

	localAPI.WebsiteUploaderHandler = operations.WebsiteUploaderHandlerFunc(WebsiteUploaderHandler)

	pluginstore.InitializePluginStoreAPI(localAPI)
	myplugin.InitializePluginAPI(localAPI)
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
		log.Fatalln(err)
	}

	localAPI := operations.NewMassastationAPI(swaggerSpec)
	server := restapi.NewServer(localAPI)

	setAPIFlags(server, flags)

	err = store.NewStore()
	if err != nil {
		log.Fatalln(err)
	}

	return &Server{
		api:      server,
		localAPI: localAPI,
		shutdown: make(chan struct{}),
	}
}

// Starts the server.
// This function starts the server in a new goroutine to avoid blocking the main thread.
func (server *Server) Start(networkManager *config.NetworkManager) {
	server.printNodeVersion(networkManager)

	initLocalAPI(server.localAPI, networkManager)
	server.api.ConfigureMassaStationAPI(*networkManager.Network(), server.shutdown)

	go func() {
		if err := server.api.Serve(); err != nil {
			log.Fatalln(err)
		}
	}()

	log.Println("Server started")
}

// Stops the server and waits for it to finish.
func (server *Server) Stop() {
	if err := server.api.Shutdown(); err != nil {
		log.Fatalln(err)
	}

	<-server.shutdown

	log.Println("Server stopped")
}

// Displays the node version of the connected node.
func (server *Server) printNodeVersion(networkManager *config.NetworkManager) {
	client := node.NewClient(networkManager.Network().NodeURL)
	status, err := node.Status(client)

	nodeVersion := "unknown"
	if err == nil {
		nodeVersion = *status.Version
	} else {
		log.Println("Could not get node version:", err)
	}

	log.Printf("Connected to node server %s (version %s)\n", networkManager.Network().NodeURL, nodeVersion)
}
