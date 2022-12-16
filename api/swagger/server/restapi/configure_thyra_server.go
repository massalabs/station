// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"embed"
	"fmt"
	"net/http"

	"github.com/go-openapi/errors"
	openApiRuntime "github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/onchain/website"
	"github.com/rs/cors"
)

//go:generate swagger generate server --target ../../server --name ThyraServer --spec ../../swagger.yml --principal interface{} --exclude-main

func configureFlags(api *operations.ThyraServerAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.ThyraServerAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = openApiRuntime.JSONConsumer()

	api.JSONProducer = openApiRuntime.JSONProducer()

	if api.CmdExecuteFunctionHandler == nil {
		api.CmdExecuteFunctionHandler = operations.CmdExecuteFunctionHandlerFunc(func(params operations.CmdExecuteFunctionParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.CmdExecuteFunctionHandler has not yet been implemented")
		})
	}

	if api.KpiHandler == nil {
		api.KpiHandler = operations.KpiHandlerFunc(func(params operations.KpiParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.Kpi has not yet been implemented")
		})
	}

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

//go:embed resource
var content embed.FS

func myCert(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	fmt.Println("inside myCert")
	certx509, err := (&mkcert{}).Run(hello.ServerName)
	var cert *tls.Certificate
	if err != nil {
		cert.Leaf = &certx509
	}

	fmt.Println("%v", hello)
	return cert, fmt.Errorf("err.Error()")
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	var err error

	var helloInfo tls.ClientHelloInfo
	helloInfo.ServerName = tlsConfig.ServerName

	tlsConfig.GetCertificate = myCert
	basePath := "resource/certificate/"

	unsecureCertificate, err := content.ReadFile(basePath + "unsecure.crt")
	if err != nil {
		panic(err)
	}

	unsecureKey, err := content.ReadFile(basePath + "unsecure.key")
	if err != nil {
		panic(err)
	}
	// Se placer la pour provide le CA une fois qu'on l'a généré
	// Récuperer le CA et ses Clefs suivant l'os via getCaROOT() et getCaKey()
	// Faire le process de certificate avec le CA et les clefs a la volée
	if len(tlsConfig.Certificates) == 0 {
		fmt.Println("warning: insecure HTTPS configuration.")
		fmt.Println("	To fix this, use your own .crt and .key files using `--tls-certificate` and `--tls-key` flags")

		var err error

		tlsConfig.Certificates = make([]tls.Certificate, 1)
		certiff, err := tls.X509KeyPair(unsecureCertificate, unsecureKey)
		tlsConfig.Certificates = append(tlsConfig.Certificates, certiff)
		if err != nil {
			fmt.Println("PANIC PANIC")
			panic(err)
		}
	}
	//Le Probleme etant que c'est appellé deux fois au lancement une premiere fois ServerName est vide
	// Une seconde fois et la ca plante donc ce tableau
	// var leafCertificate *x509.Certificate

	// Ok donc mon certificat est donc dans leafCertificate maintenant on fait quoi
	if err != nil {
		fmt.Println("C CASSE CHEF")
	}
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	handleCORS := cors.Default().Handler

	return website.TopMiddleware(handleCORS(handler))
}
