package api

import (
	"log"
	"os"
	"sync"

	"fyne.io/fyne/v2"
	"github.com/go-openapi/loads"
	"github.com/jessevdk/go-flags"
	"github.com/massalabs/thyra/api/swagger/server/restapi"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/int/api/cmd"
	"github.com/massalabs/thyra/int/api/myplugin"
	"github.com/massalabs/thyra/int/api/plugin"
	"github.com/massalabs/thyra/int/api/wallet"
	"github.com/massalabs/thyra/int/api/websites"
	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/onchain/dns"
	pluginmanager "github.com/massalabs/thyra/pkg/plugins"
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

	parseNetworkFlag(&startFlags.MassaNodeServer)

	if startFlags.DNSAddress != "" {
		os.Setenv(dns.EnvKey, startFlags.DNSAddress)
	}
}

func parseNetworkFlag(massaNodeServerPtr *string) {
	var dnsAddress string

	switch *massaNodeServerPtr {
	case "TESTNET":
		*massaNodeServerPtr = "https://test.massa.net/api/v2"
		// testnet18
		dnsAddress = "A13Ten7uRWQrLHYfve4NXg13aiWxMYWMPM2dWfbVBEB2w7eYC1J"
	case "LABNET":
		*massaNodeServerPtr = "https://labnet.massa.net/api/v2"
		dnsAddress = "A12RgLPuRQaVTue2CtPws6deXUfUnk6nfveZS9bedyzoNS8WyYtg"

	case "INNONET":
		*massaNodeServerPtr = "https://inno.massa.net/test17"
		dnsAddress = "A12WHCz3Qf4iPSaWMWA1ErbCUxJrHQSfY24BrPte3HELsCG7YxJh"

	case "LOCALHOST":
		*massaNodeServerPtr = "http://127.0.0.1:33035"
	}

	os.Setenv(dns.EnvKey, dnsAddress)
	os.Setenv("MASSA_NODE_URL", *massaNodeServerPtr)
}

func stopServer(app *fyne.App, server *restapi.Server, manager *pluginmanager.PluginManager) {
	manager.StopPlugins()

	if err := server.Shutdown(); err != nil {
		log.Fatalln(err)
	}

	(*app).Quit()
}

func initLocalAPI(localAPI *operations.ThyraServerAPI, app *fyne.App, manager *pluginmanager.PluginManager) {
	var walletStorage sync.Map

	localAPI.CmdExecuteFunctionHandler = operations.CmdExecuteFunctionHandlerFunc(
		cmd.CreateExecuteFunctionHandler(app))

	localAPI.MgmtPluginsListHandler = plugin.NewGet(manager)

	localAPI.MgmtWalletGetHandler = wallet.NewGet(&walletStorage)
	localAPI.MgmtWalletCreateHandler = wallet.NewCreate(&walletStorage)
	localAPI.MgmtWalletImportHandler = wallet.NewImport(&walletStorage, app)
	localAPI.MgmtWalletDeleteHandler = wallet.NewDelete(&walletStorage, app)

	localAPI.WebsiteCreatorPrepareHandler = operations.WebsiteCreatorPrepareHandlerFunc(
		websites.CreatePrepareForWebsiteHandler(app),
	)
	localAPI.WebsiteCreatorUploadHandler = operations.WebsiteCreatorUploadHandlerFunc(
		websites.CreateUploadWebsiteHandler(app),
	)
	localAPI.WebsiteUploadMissingChunksHandler = operations.WebsiteUploadMissingChunksHandlerFunc(
		websites.CreateUploadMissingChunksHandler(app),
	)
	localAPI.MyDomainsGetterHandler = operations.MyDomainsGetterHandlerFunc(websites.DomainsHandler)
	localAPI.AllDomainsGetterHandler = operations.AllDomainsGetterHandlerFunc(websites.RegistryHandler)

	localAPI.ThyraRegistryHandler = operations.ThyraRegistryHandlerFunc(ThyraRegistryHandler)

	localAPI.ThyraEventsGetterHandler = operations.ThyraEventsGetterHandlerFunc(EventListenerHandler)
	localAPI.BrowseHandler = operations.BrowseHandlerFunc(BrowseHandler)

	localAPI.ThyraWalletHandler = operations.ThyraWalletHandlerFunc(ThyraWalletHandler)
	localAPI.ThyraWebsiteCreatorHandler = operations.ThyraWebsiteCreatorHandlerFunc(ThyraWebsiteCreatorHandler)

	myplugin.InitializePluginAPI(localAPI)
}

func StartServer(app *fyne.App, startFlags StartServerFlags) {
	// Initialize Swagger
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	localAPI := operations.NewThyraServerAPI(swaggerSpec)
	server := restapi.NewServer(localAPI)

	setAPIFlags(server, startFlags)

	// Display info about node server
	client := node.NewDefaultClient()
	status, err := node.Status(client)

	nodeVersion := "unknown"
	if err == nil {
		nodeVersion = *status.Version
	} else {
		log.Println("Could not get node version:", err)
	}

	log.Printf("Connected to node server %s (version %s)\n", os.Getenv("MASSA_NODE_URL"), nodeVersion)

	// Load plugins
	manager, err := pluginmanager.New(server.Port, server.TLSPort)
	if err != nil {
		log.Fatalln(err)
	}

	// Start plugins
	manager.StartPlugins()
	defer stopServer(app, server, manager)

	initLocalAPI(localAPI, app, manager)
	server.ConfigureAPI()

	if err := server.Serve(); err != nil {
		//nolint:gocritic
		log.Fatalln(err)
	}
}
