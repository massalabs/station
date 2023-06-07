package api

import (
	"log"
	"net/http"

	"github.com/massalabs/thyra/api/interceptor"
	"github.com/massalabs/thyra/pkg/config"
	"github.com/massalabs/thyra/pkg/onchain/website"
	"github.com/massalabs/thyra/pkg/plugin"
)

// TopMiddleware is called by go-swagger framework before its endpoints.
// current defined interceptor are:
// - MassaTLDInterceptor to handle *.massa websites
// - Plugin interceptor to handle call to registered plugins
// - Default resource interceptor to handle browser call (needed for mobile?) and web resources not yet pluginized.
func TopMiddleware(handler http.Handler, config config.AppConfig) http.Handler {
	//nolint:varnamelen
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s %s]", r.Method, r.URL.Path)
		// Goes through all local interceptors.
		req := website.RedirectToDefaultResourceInterceptor(
			plugin.Interceptor(
				website.MassaTLDInterceptor(&interceptor.Interceptor{Writer: w, Request: r}, config))) //nolint:contextcheck
		// if the request was not handled by any interceptor, let the swagger API takes care of it.
		if req != nil {
			handler.ServeHTTP(w, r)
		}
	})
}
