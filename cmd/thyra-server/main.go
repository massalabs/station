package main

import (
	"log"

	"github.com/go-openapi/loads"
	"github.com/massalabs/thyra/api/swagger/server/restapi"

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
	server.TLSCertificate = "insecure.crt"
	server.TLSCertificateKey = "insecure.key"

	// Start server which listening
	server.ConfigureAPI()
	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}
