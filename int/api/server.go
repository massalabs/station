package api

import (
	_ "embed"
	"flag"
	"log"
	"os"
	"sync"

	"fyne.io/fyne/v2"
	"github.com/go-openapi/loads"
	"github.com/jessevdk/go-flags"
	"github.com/massalabs/thyra/api/swagger/server/restapi"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/int/api/cmd"
	"github.com/massalabs/thyra/int/api/plugin"
	"github.com/massalabs/thyra/int/api/wallet"
	"github.com/massalabs/thyra/int/api/websites"
	pluginmanager "github.com/massalabs/thyra/pkg/plugins"
)

//go:embed version.txt
var version string

func parseFlags(server *restapi.Server) bool {
	const httpPort = 80

	const httpsPort = 443

	flag.IntVar(&server.Port, "http-port", httpPort, "HTTP port to listen to")

	flag.IntVar(&server.TLSPort, "https-port", httpsPort, "HTTPS port to listen to")

	certFilePtr := flag.String("tls-certificate", "", "path to certificate file")
	keyFilePtr := flag.String("tls-key", "", "path to key file")
	massaNodeServerPtr := flag.String("node-server", "TESTNET", `Massa node that Thyra connects to. 
	Can be an IP address, a URL or one of the following values: 'TESTNET', 'LABNET', 'INNONET' or LOCALHOST`)
	version := flag.Bool("version", false, "Print version info and exit")

	flag.Parse()

	if *certFilePtr != "" {
		server.TLSCertificate = flags.Filename(*certFilePtr)
	}

	if *keyFilePtr != "" {
		server.TLSCertificateKey = flags.Filename(*keyFilePtr)
	}

	parseNetworkFlag(massaNodeServerPtr)

	return *version
}

func parseNetworkFlag(massaNodeServerPtr *string) {
	switch *massaNodeServerPtr {
	case "TESTNET":
		*massaNodeServerPtr = "https://test.massa.net/api/v2"
	case "LABNET":
		*massaNodeServerPtr = "https://labnet.massa.net/"
	case "INNONET":
		*massaNodeServerPtr = "https://inno.massa.net/test15"
	case "LOCALHOST":
		*massaNodeServerPtr = "http://127.0.0.1:33035"
	}

	os.Setenv("MASSA_NODE_URL", *massaNodeServerPtr)
}

func stopServer(app *fyne.App, server *restapi.Server, manager *pluginmanager.PluginManager) {
	defer func() {
		manager.StopPlugins()

		if err := server.Shutdown(); err != nil {
			log.Fatalln(err)
		}

		(*app).Quit()
	}()
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
}

func StartServer(app *fyne.App) {
	// Initialize Swagger
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	localAPI := operations.NewThyraServerAPI(swaggerSpec)
	server := restapi.NewServer(localAPI)

	hasVersionFlag := parseFlags(server)

	// Run plugins
	manager, err := pluginmanager.New(server.Port, server.TLSPort)
	if err != nil {
		log.Fatalln(err)
	}

	if hasVersionFlag {
		log.Println("Thyra version", version)
		stopServer(app, server, manager)

		return
	}
	initLocalAPI(localAPI, app, manager)
	server.ConfigureAPI()

	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}

	stopServer(app, server, manager)
}
