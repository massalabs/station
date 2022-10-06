package main

import (
	"flag"
	"log"
	"os"
	"sync"

	"github.com/go-openapi/loads"
	"github.com/jessevdk/go-flags"
	"github.com/massalabs/thyra/api/swagger/server/restapi"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/int/api"
	"github.com/massalabs/thyra/int/api/cmd"
	"github.com/massalabs/thyra/int/api/wallet"
	"github.com/massalabs/thyra/int/api/websites"
)

func parseFlags(server *restapi.Server) {
	const httpPort = 80

	const httpsPort = 443

	flag.IntVar(&server.Port, "http-port", httpPort, "HTTP port to listen to")

	flag.IntVar(&server.TLSPort, "https-port", httpsPort, "HTTPS port to listen to")

	certFilePtr := flag.String("tls-certificate", "", "path to certificate file")
	keyFilePtr := flag.String("tls-key", "", "path to key file")
	massaNodeServerPtr := flag.String("node-server", "INNONET", `Massa node that Thyra connects to. 
	Can be an IP address, a URL or one of the following values: 'TESTNET', 'LABNET', 'INNONET' or LOCALHOST`)

	flag.Parse()

	if *certFilePtr != "" {
		server.TLSCertificate = flags.Filename(*certFilePtr)
	}

	if *keyFilePtr != "" {
		server.TLSCertificateKey = flags.Filename(*keyFilePtr)
	}

	parseNetworkFlag(massaNodeServerPtr)
}

func parseNetworkFlag(massaNodeServerPtr *string) {
	switch *massaNodeServerPtr {
	case "TESTNET":
		*massaNodeServerPtr = "https://test.massa.net/v1/"
	case "LABNET":
		*massaNodeServerPtr = "https://labnet.massa.net/"
	case "INNONET":
		*massaNodeServerPtr = "https://inno.massa.net/test13"
	case "LOCALHOST":
		*massaNodeServerPtr = "http://127.0.0.1"
	}

	os.Setenv("MASSA_NODE_URL", *massaNodeServerPtr)
}

func main() {
	// Initialize Swagger
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	localAPI := operations.NewThyraServerAPI(swaggerSpec)
	server := restapi.NewServer(localAPI)

	defer func() {
		if err := server.Shutdown(); err != nil {
			log.Fatalln(err)
		}
	}()

	parseFlags(server)

	var walletStorage sync.Map

	localAPI.CmdExecuteFunctionHandler = operations.CmdExecuteFunctionHandlerFunc(cmd.ExecuteFunctionHandler)

	localAPI.MgmtWalletGetHandler = wallet.NewGet(&walletStorage)
	localAPI.MgmtWalletCreateHandler = wallet.NewCreate(&walletStorage)
	localAPI.MgmtWalletImportHandler = wallet.NewImport(&walletStorage)
	localAPI.MgmtWalletDeleteHandler = wallet.NewDelete(&walletStorage)

	localAPI.WebsiteCreatorPrepareHandler = operations.WebsiteCreatorPrepareHandlerFunc(websites.PrepareForWebsiteHandler)
	localAPI.WebsiteCreatorUploadHandler = operations.WebsiteCreatorUploadHandlerFunc(websites.UploadWebsiteHandler)
	localAPI.MyDomainsGetterHandler = operations.MyDomainsGetterHandlerFunc(websites.DomainsHandler)
	localAPI.AllDomainsGetterHandler = operations.AllDomainsGetterHandlerFunc(websites.RegistryHandler)

	localAPI.ThyraRegistryHandler = operations.ThyraRegistryHandlerFunc(api.ThyraRegistryHandler)

	localAPI.ThyraEventsGetterHandler = operations.ThyraEventsGetterHandlerFunc(api.EventListenerHandler)
	localAPI.BrowseHandler = operations.BrowseHandlerFunc(api.BrowseHandler)

	localAPI.ThyraWalletHandler = operations.ThyraWalletHandlerFunc(api.ThyraWalletHandler)
	localAPI.ThyraWebsiteCreatorHandler = operations.ThyraWebsiteCreatorHandlerFunc(api.ThyraWebsiteCreatorHandler)

	server.ConfigureAPI()

	if err := server.Serve(); err != nil {
		//nolint:gocritic
		log.Fatalln(err)
	}
}
