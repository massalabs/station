package main

import (
	"flag"
	"log"

	"github.com/go-openapi/loads"
	"github.com/jessevdk/go-flags"
	"github.com/massalabs/thyra/api/swagger/server/restapi"
	apiHandler "github.com/massalabs/thyra/internal/apihandler"

	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
)

func main() {

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

	api.CmdExecuteFunctionHandler = operations.CmdExecuteFunctionHandlerFunc(apiHandler.ExecuteFunction)

	// Start server which listening
	server.ConfigureAPI()
	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}
