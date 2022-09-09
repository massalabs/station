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
)

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

	flag.IntVar(&server.Port, "http-port", 80, "HTTP port to listen to")

	flag.IntVar(&server.TLSPort, "https-port", 443, "HTTPS port to listen to")

	certFilePtr := flag.String("tls-certificate", "", "path to certificate file")
	keyFilePtr := flag.String("tls-key", "", "path to key file")
	massaNetworkAddressPtr := flag.String("network", "", "IP address of the chosen massa network")
	flag.Parse()

	if *certFilePtr != "" {
		server.TLSCertificate = flags.Filename(*certFilePtr)
	}

	if *keyFilePtr != "" {
		server.TLSCertificateKey = flags.Filename(*keyFilePtr)
	}

	if *massaNetworkAddressPtr != "" {
		if *massaNetworkAddressPtr == "TESTNET" {
			*massaNetworkAddressPtr = "https://test.massa.net/v1/"
		} else if *massaNetworkAddressPtr == "LABNET" {
			*massaNetworkAddressPtr = "https://labnet.massa.net/"
		} else if *massaNetworkAddressPtr == "INNONET" {
			*massaNetworkAddressPtr = "https://inno.massa.net/test13"
		} else if *massaNetworkAddressPtr == "LOCALHOST" {
			*massaNetworkAddressPtr = "http://127.0.0.1"
		}
		os.Setenv("MASSA_URL_RPC", *massaNetworkAddressPtr)
	}

	var walletStorage sync.Map

	localAPI.CmdExecuteFunctionHandler = cmd.NewExecuteFunction(&walletStorage)

	localAPI.MgmtWalletGetHandler = wallet.NewGet(&walletStorage)
	localAPI.MgmtWalletCreateHandler = wallet.NewCreate(&walletStorage)
	localAPI.MgmtWalletImportHandler = wallet.NewImport(&walletStorage)
	localAPI.MgmtWalletDeleteHandler = wallet.NewDelete(&walletStorage)

	localAPI.WebsiteCreatorPrepareHandler = operations.WebsiteCreatorPrepareHandlerFunc(api.PrepareForWebsiteHandler)
	localAPI.WebsiteCreatorUploadHandler = operations.WebsiteCreatorUploadHandlerFunc(api.UploadWebsiteHandler)

	localAPI.MyDomainsGetterHandler = operations.MyDomainsGetterHandlerFunc(api.DomainsHandler)

	localAPI.BrowseHandler = operations.BrowseHandlerFunc(api.BrowseHandler)

	localAPI.ThyraWalletHandler = operations.ThyraWalletHandlerFunc(api.ThyraWalletHandler)
	localAPI.ThyraWebsiteCreatorHandler = operations.ThyraWebsiteCreatorHandlerFunc(api.ThyraWebsiteCreatorHandler)

	// Start server which listening
	server.ConfigureAPI()
	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}
