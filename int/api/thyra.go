package api

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/front"
	"github.com/massalabs/thyra/pkg/front/wallet"
	"github.com/massalabs/thyra/pkg/front/website"
)

func ThyraWalletHandler(params operations.ThyraWalletParams) middleware.Responder {
	var body string

	switch params.Resource {
	case "wallet.css":
		body = wallet.CSS
	case "wallet.js":
		body = wallet.JS
	case "index.html":
		body = wallet.HTML
	case "logo_banner.webp":
		body = front.LogoBanner
	case "logo.png":
		body = front.Logo
	case "errors.js":
		body = front.Errors

	}

	return NewCustomResponder([]byte(body), contentType(params.Resource), http.StatusOK)
}

func ThyraWebsiteCreatorHandler(params operations.ThyraWebsiteCreatorParams) middleware.Responder {
	var body string

	switch params.Resource {
	case "website.css":
		body = website.CSS
	case "website.js":
		body = website.JS
	case "index.html":
		body = website.HTML
	case "logo_banner.webp":
		body = front.LogoBanner
	case "logo.png":
		body = front.Logo
	case "errors.js":
		body = front.Errors

	}

	return NewCustomResponder([]byte(body), contentType(params.Resource), http.StatusOK)
}
