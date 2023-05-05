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

func stopServer(server *restapi.Server) {
	if err := server.Shutdown(); err != nil {
		log.Fatalln(err)
	}
}

func initLocalAPI(localAPI *operations.ThyraServerAPI, config config.AppConfig) {
	localAPI.CmdExecuteFunctionHandler = cmd.NewExecuteFunction(&config)

	localAPI.MassaGetAddressesHandler = massa.NewGetAddress(&config)

	localAPI.CmdDeploySCHandler = cmd.NewDeploySC(&config)

	localAPI.WebsiteCreatorPrepareHandler = websites.NewWebsitePrepareHandler(&config)
	localAPI.WebsiteCreatorUploadHandler = websites.NewWebsiteUploadHandler(&config)
	localAPI.WebsiteUploadMissingChunksHandler = websites.NewWebsiteUploadMissedChunkHandler(&config)

	localAPI.MyDomainsGetterHandler = websites.NewDomainsHandler(&config)
	localAPI.AllDomainsGetterHandler = websites.NewRegistryHandler(&config)

	localAPI.ThyraRegistryHandler = operations.ThyraRegistryHandlerFunc(ThyraRegistryHandler)
	localAPI.ThyraHomeHandler = operations.ThyraHomeHandlerFunc(ThyraHomeHandler)
	localAPI.ThyraEventsGetterHandler = NewEventListenerHandler(&config)
	localAPI.BrowseHandler = NewBrowseHandler(&config)
	localAPI.ThyraPluginManagerHandler = operations.ThyraPluginManagerHandlerFunc(ThyraPluginManagerHandler)

	localAPI.ThyraWebsiteCreatorHandler = operations.ThyraWebsiteCreatorHandlerFunc(ThyraWebsiteCreatorHandler)

	myplugin.InitializePluginAPI(localAPI)
	pluginstore.InitializePluginStoreAPI(localAPI)
}

func StartServer(flags StartServerFlags) {
	// Initialize Swagger
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	localAPI := operations.NewThyraServerAPI(swaggerSpec)
	server := restapi.NewServer(localAPI)

	setAPIFlags(server, flags)

	config := config.AppConfig{
		NodeURL:    config.GetNodeURL(flags.MassaNodeServer),
		DNSAddress: config.GetDNSAddress(flags.MassaNodeServer, flags.DNSAddress),
		Network:    config.GetNetwork(flags.MassaNodeServer),
	}

	// Display info about node server
	client := node.NewClient(config.NodeURL)
	status, err := node.Status(client)

	nodeVersion := "unknown"
	if err == nil {
		nodeVersion = *status.Version
	} else {
		log.Println("Could not get node version:", err)
	}

	log.Printf("Connected to node server %s (version %s)\n", config.NodeURL, nodeVersion)

	defer stopServer(server)

	initLocalAPI(localAPI, config)

	server.ConfigureMassaStationAPI(config)

	if err := server.Serve(); err != nil {
		//nolint:gocritic
		log.Fatalln(err)
	}
}
