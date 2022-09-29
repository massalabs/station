package api

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/front"
	"github.com/massalabs/thyra/pkg/front/registry"
	"github.com/massalabs/thyra/pkg/front/wallet"
	"github.com/massalabs/thyra/pkg/front/website"
)

const (
	indexHTML  = "index.html"
	logoBanner = "logo_banner.webp"
	logoPNG    = "logo.png"
	errorsJS   = "errors.js"
	commonJS   = "common.js"
)

//nolint:nolintlint,ireturn
func ThyraWalletHandler(params operations.ThyraWalletParams) middleware.Responder {
	var body string

	switch params.Resource {
	case "wallet.css":
		body = wallet.CSS
	case "wallet.js":
		body = wallet.JS
	case indexHTML:
		body = wallet.HTML
	case logoBanner:
		body = front.LogoBanner
	case logoPNG:
		body = front.Logo
	case errorsJS:
		body = front.Errors
	case commonJS:
		body = front.Common
	}

	return NewCustomResponder([]byte(body), contentType(params.Resource), http.StatusOK)
}

//nolint:nolintlint,ireturn
func ThyraWebsiteCreatorHandler(params operations.ThyraWebsiteCreatorParams) middleware.Responder {
	var body string

	switch params.Resource {
	case "website.css":
		body = website.CSS
	case "website.js":
		body = website.JS
	case indexHTML:
		body = website.HTML
	case logoBanner:
		body = front.LogoBanner
	case logoPNG:
		body = front.Logo
	case errorsJS:
		body = front.Errors
	case commonJS:
		body = front.Common
	case "event-manager.js":
		body = front.EventListener
	}

	return NewCustomResponder([]byte(body), contentType(params.Resource), http.StatusOK)
}

//nolint:nolintlint,ireturn
func ThyraRegistryHandler(params operations.ThyraRegistryParams) middleware.Responder {
	var body string

	switch params.Resource {
	case indexHTML:
		body = registry.HTML
	case "registry.js":
		body = registry.JS
	case "registry.css":
		body = registry.CSS
	case logoBanner:
		body = front.LogoBanner
	case logoPNG:
		body = front.Logo
	case errorsJS:
		body = front.Errors
	case commonJS:
		body = front.Common
	case "event-manager.js":
		body = front.EventListener
	}

	return NewCustomResponder([]byte(body), contentType(params.Resource), http.StatusOK)
}
