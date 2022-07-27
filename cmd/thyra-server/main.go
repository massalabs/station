package main

import (
	"flag"
	"log"
	"sync"

	"github.com/go-openapi/loads"
	"github.com/jessevdk/go-flags"
	"github.com/massalabs/thyra/api/swagger/server/restapi"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/int/apihandler/cmd"
	"github.com/massalabs/thyra/int/apihandler/wallet"
)

//TODO Manage file generation in an other way, generateFiles is not working while deployed
func main() {
	// Generate files
	//front.GenerateFiles()
	// Initialize Swagger
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewThyraServerAPI(swaggerSpec)
	server := restapi.NewServer(api)

	defer func() {
		if err := server.Shutdown(); err != nil {
			log.Fatalln(err)
		}
	}()

	server.Port = 80
	server.TLSPort = 443

	certFilePtr := flag.String("tls-certificate", "", "path to certificate file")
	keyFilePtr := flag.String("tls-key", "", "path to key file")
	flag.Parse()

	if *certFilePtr != "" {
		server.TLSCertificate = flags.Filename(*certFilePtr)
	}

	if *keyFilePtr != "" {
		server.TLSCertificateKey = flags.Filename(*keyFilePtr)
	}

	var walletStorage sync.Map

	api.CmdExecuteFunctionHandler = cmd.NewExecuteFunction(&walletStorage)

	api.MgmtWalletGetHandler = wallet.NewGet(&walletStorage)
	api.MgmtWalletCreateHandler = wallet.NewCreate(&walletStorage)
	api.MgmtWalletImportHandler = wallet.NewImport(&walletStorage)
	api.MgmtWalletDeleteHandler = wallet.NewDelete(&walletStorage)

	// Start server which listening
	server.ConfigureAPI()
	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}
