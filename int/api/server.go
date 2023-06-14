package api

import (
	"log"
	"os"

	"github.com/go-openapi/loads"
	"github.com/jessevdk/go-flags"
	"github.com/massalabs/thyra/api/swagger/server/restapi"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/int/api/cmd"
	"github.com/massalabs/thyra/int/api/massa"
	"github.com/massalabs/thyra/int/api/myplugin"
	"github.com/massalabs/thyra/int/api/pluginstore"
	"github.com/massalabs/thyra/int/api/websites"
	"github.com/massalabs/thyra/pkg/config"
	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/store"
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

func initLocalAPI(localAPI *operations.MassastationAPI, config config.AppConfig) {
	localAPI.CmdExecuteFunctionHandler = cmd.NewExecuteFunctionHandler(&config)

	localAPI.MassaGetAddressesHandler = massa.NewGetAddressHandler(&config)
	localAPI.GetNodeHandler = massa.NewGetNodeHandler(&config)

	localAPI.CmdDeploySCHandler = cmd.NewDeploySCHandler(&config)

	localAPI.WebsiteUploaderPrepareHandler = websites.NewWebsitePrepareHandler(&config)
	localAPI.WebsiteUploaderUploadHandler = websites.NewWebsiteUploadHandler(&config)
	localAPI.WebsiteUploadMissingChunksHandler = websites.NewWebsiteUploadMissedChunkHandler(&config)

	localAPI.MyDomainsGetterHandler = websites.NewDomainsHandler(&config)
	localAPI.AllDomainsGetterHandler = websites.NewRegistryHandler(&config)

	localAPI.WebOnChainSearchHandler = operations.WebOnChainSearchHandlerFunc(WebOnChainSearchHandler)
	localAPI.MassaStationHomeHandler = operations.MassaStationHomeHandlerFunc(MassaStationHomeHandler)
	localAPI.EventsGetterHandler = NewEventListenerHandler(&config)
	localAPI.MassaStationPluginManagerHandler = operations.MassaStationPluginManagerHandlerFunc(MassaStationPluginManagerHandler)
	localAPI.MassaStationWebAppHandler = operations.MassaStationWebAppHandlerFunc(MassaStationWebAppHandler)

	localAPI.WebsiteUploaderHandler = operations.WebsiteUploaderHandlerFunc(WebsiteUploaderHandler)
	pluginstore.InitializePluginStoreAPI(localAPI)
	myplugin.InitializePluginAPI(localAPI)
}

type Server struct {
	config   config.AppConfig
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

	config := config.AppConfig{
		NodeURL:    config.GetNodeURL(flags.MassaNodeServer),
		DNSAddress: config.GetDNSAddress(flags.MassaNodeServer, flags.DNSAddress),
		Network:    config.GetNetwork(flags.MassaNodeServer),
	}

	return &Server{
		config:   config,
		api:      server,
		localAPI: localAPI,
		shutdown: make(chan struct{}),
	}
}

// Starts the server.
// This function starts the server in a new goroutine to avoid blocking the main thread.
func (server *Server) Start() {
	server.printNodeVersion()

	initLocalAPI(server.localAPI, server.config)
	server.api.ConfigureMassaStationAPI(server.config, server.shutdown)

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
func (server *Server) printNodeVersion() {
	client := node.NewClient(server.config.NodeURL)
	status, err := node.Status(client)

	nodeVersion := "unknown"
	if err == nil {
		nodeVersion = *status.Version
	} else {
		log.Println("Could not get node version:", err)
	}

	log.Printf("Connected to node server %s (version %s)\n", server.config.NodeURL, nodeVersion)
}
