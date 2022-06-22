// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

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

	if api.CmdCallSCHandler == nil {
		api.CmdCallSCHandler = operations.CmdCallSCHandlerFunc(func(params operations.CmdCallSCParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.CmdCallSC has not yet been implemented")
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

		certPem := []byte(`-----BEGIN CERTIFICATE-----
MIIEqzCCApMCCQDSu9BWWGogFjANBgkqhkiG9w0BAQsFADAVMRMwEQYDVQQDDApH
byBTd2FnZ2VyMB4XDTE2MTIyMTA4NDEzOVoXDTE3MTIyMTA4NDEzOVowGjEYMBYG
A1UEAwwPZ29zd2FnZ2VyLmxvY2FsMIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIIC
CgKCAgEA55Nv8PnHB6HUTknFqWWrQ7a8na3P3SxBOgszFok32vaK78c5mSPwgg2b
zb6Mi0yLCYBmQGnF2EGMtmTzjCEYiZYIWHNNN3oo4xm108SLzP+/E7swTGgJ+zcK
O+eHWEKtdFipWNIRvpPqmn7Pm8dA6dPEDi8uqc1GcRprtZayWY3A++M3dptqEvbw
uz0/utVE3c/1Q4ltraZeasjNxIzR0gdRG5Ty/wG8k5L3pLY/oTZ0ZExpHOVNa0WF
9OgMKlZMnsir/8jgO5j30eUrtFkiQ3/tM5+jX7k+TMDVk0I1XiIGas+6Ja2GNhKp
yrhZeom9PTHtyW1KVS43u82hmh5Cedz+IbtmM6ZIwtmZFM2lg7fKkGgqJm3acEsH
LkK9fpALmBdxiSoRYrlzB9Z/DIJ23iNReX2vOmvdvLc5tdgflchq+qR6NOok4zcw
QC7imgjPXRnvdFiyRKYvXwWuC745MKFpcpUWXGiRwmiVTWBfEzshd1gL2tMcNvZW
eTHRAMGGbJelSlsnRHsk4eQq4W29uJ1ijqRQKavqJNr14pMXE8a/0KiLmlU22CMe
GL/+9JHudNa9Q8sSu4Ils+GPKgmjI6llMXq9v5VHD2b7OBX5XCbfkldP8zNNXMeR
NTAYQixyVlqihf7qsUQGFbQ1jkkh60h7RAZJldV+QK9mjX/+k5UCAwEAATANBgkq
hkiG9w0BAQsFAAOCAgEAU4qxBUGC3HNsL40nf7LEJ4P5li+/93jPBtNKiwkMW0tJ
m69uuM8sFniyju9v+ity/voOLWScZuhF07YGThXriaOOUrFrSfV2cWlv95kU+oJA
ux4xgjZX4OXAkm3iZJyaHOeWh6zMjt0SFdJfsYydDP5/qw/4aWsRon/Szv8lDE8W
fd0nYqbz5YgwZiR0yF7VJYp2UMStEg0kINL1rQR7DvWTba3mkM+QaOpoP1HGCOCZ
GCaREsbXsVp3vfPhI480Uxixj49LNjZZfa7sWvHgBjeQBckK47m90JxqPJrF5yfz
zKyAMBxzhU8pWAptaTdhC6OKC01lUe5CQDiFwlKATk5/aliQNLmaVPF4x87aoWUl
oUBvODz/jDBPWeLH5IhdYMcmCQne0C+dn2R7M5OBqtJaasIPBkq0iqQYWs6ww/jI
Ugo2WeUXl6fqgZGn/Q1yR7ITPV5bskoh7ObikPjLWTlQMm906lV7vcL+R/3AOc8X
Yj5pwQhRoAmpFCoN90nLqAd1tEHhl8K6ILoFk6cHTKQNMgR2veuGGWouy1h9Z3S+
BendwOhHed10tPVbezvZlj2avzc6lESLDlgncp//KtGwNMKoXfBI5bhR34uKVgm7
pYuOPuO7zzpEhNlDlmhos8jHZIbsE8dFPNL9WCblLZMSYFXfLqH0984IcDhM26M=
-----END CERTIFICATE-----`)
		keyPem := []byte(`-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`)

		var err error

		tlsConfig.Certificates = make([]tls.Certificate, 1)
		tlsConfig.Certificates[0], err = tls.X509KeyPair(certPem, keyPem)
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
	return website.HandlerFunc(handler)
}
