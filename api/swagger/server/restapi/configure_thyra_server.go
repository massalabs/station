// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/massalabs/thyra/api/swagger/server"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/onchain/website"
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

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

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

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	if len(tlsConfig.Certificates) == 0 {
		fmt.Println("warning: insecure HTTPS configuration.")
		fmt.Println("	To fix this, use your own .crt and .key files using `--tls-certificate` and `--tls-key` flags")

		var err error

		tlsConfig.Certificates = make([]tls.Certificate, 1)
		tlsConfig.Certificates[0], err = tls.X509KeyPair([]byte(server.UnsecureCertificate), []byte(server.UnsecureKey))

		if err != nil {
			panic(err)
		}
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
	return website.TopMiddleware(handler)
}
