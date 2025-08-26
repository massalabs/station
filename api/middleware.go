package api

import (
	"net/http"

	"github.com/massalabs/station/api/interceptor"
	"github.com/massalabs/station/pkg/logger"
	"github.com/massalabs/station/pkg/plugin"
)

// TopMiddleware is called by go-swagger framework before its endpoints.
// current defined interceptor are:
// - Plugin interceptor to handle call to registered plugins
// - Default resource interceptor to handle browser call (needed for mobile?) and web resources not yet pluginized.
func TopMiddleware(handler http.Handler) http.Handler {
	// Apply domain restriction middleware first
	handler = DomainRestrictionMiddleware(handler)

	//nolint:varnamelen
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Infof("[%s %s]", r.Method, r.URL.Path)

		// Goes through all local interceptors.
		req := RedirectToDefaultResourceInterceptor(
			plugin.Interceptor(&interceptor.Interceptor{Writer: w, Request: r}))
		// if the request was not handled by any interceptor, let the swagger API takes care of it.
		if req != nil {
			handler.ServeHTTP(w, r)
		}
	})
}

// RedirectToDefaultResourceInterceptor redirects to default page or adds /index.html if the request must be handled
// by a web resource handler: routes ending with /{resource} for front-ends.
func RedirectToDefaultResourceInterceptor(req *interceptor.Interceptor) *interceptor.Interceptor {
	if req == nil {
		return nil
	}

	// redirect / to /home/index.html
	if req.Request.URL.Path == "/" {
		http.Redirect(
			req.Writer,
			req.Request,
			"/web/index.html",
			http.StatusSeeOther,
		)

		return nil
	}

	redirectPaths := []string{"/web", "/home", "/search", "/store"}

	for _, path := range redirectPaths {
		if req.Request.URL.Path == path || req.Request.URL.Path == path+"/" {
			http.Redirect(
				req.Writer,
				req.Request,
				path+"/index.html",
				http.StatusSeeOther,
			)

			return nil
		}
	}

	return req
}
