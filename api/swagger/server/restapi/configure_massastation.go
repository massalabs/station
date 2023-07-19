// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"embed"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station/api"
	"github.com/massalabs/station/api/swagger/server/restapi/operations"
	"github.com/massalabs/station/int/sni"
	MSConfig "github.com/massalabs/station/pkg/config"
	"github.com/massalabs/station/pkg/logger"
	"github.com/rs/cors"
)

func configureFlags(api *operations.MassastationAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.MassastationAPI) http.Handler {
	// unused
	return nil
}

func (s *Server) ConfigureMassaStationAPI(config MSConfig.AppConfig, shutdown chan struct{}) {
	if s.api != nil {
		s.handler = configureMassaStationAPI(s.api, config, shutdown)
	}
}

func configureMassaStationAPI(api *operations.MassastationAPI, config MSConfig.AppConfig, shutdown chan struct{}) http.Handler {
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
		api.CmdExecuteFunctionHandler = operations.CmdExecuteFunctionHandlerFunc(
			func(params operations.CmdExecuteFunctionParams) middleware.Responder {
				return middleware.NotImplemented("operation operations.CmdExecuteFunctionHandler has not yet been implemented")
			})
	}

	if api.KpiHandler == nil {
		api.KpiHandler = operations.KpiHandlerFunc(func(params operations.KpiParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.Kpi has not yet been implemented")
		})
	}

	api.PreServerShutdown = func() {
		// Create a timeout channel to force the execution of s.api.ServerShutdown()
		timeout := time.After(5 * time.Second)
		// Use a separate goroutine to handle the forced execution
		go func() {
			<-timeout
			api.ServerShutdown()
		}()
	}

	api.ServerShutdown = func() {
		close(shutdown)
	}
	api.Logger = func(msg string, args ...interface{}) {
		logger.Infof(msg, args...)
	}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares), config)
}

//go:embed resource
var content embed.FS

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	tlsConfig.GetCertificate = sni.GenerateTLS
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(_ *http.Server, _, _ string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json
// document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler, config MSConfig.AppConfig) http.Handler {
	handleCORS := cors.Default().Handler

	return api.TopMiddleware(handleCORS(handler), config)
}
