package website

import (
	"net/http"
	"strings"

	"github.com/massalabs/thyra/api/interceptor"
	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/onchain/dns"
)

func handleMassaDomainRequest(writer http.ResponseWriter, reader *http.Request, index int) {
	name := reader.Host[:index]

	rpcClient := node.NewDefaultClient()

	addr, err := dns.Resolve(rpcClient, name)
	if err != nil {
		panic(err)
	}

	var target string
	if reader.URL.Path == "/" {
		target = "index.html"
	} else {
		target = reader.URL.Path[1:]
	}

	Request(writer, reader, rpcClient, addr, target)
}

// Redirects a given request to https.
func redirectTLDToHTTPS(req *interceptor.Interceptor) {
	url := "https://" + req.Request.Host + req.Request.URL.Path

	if len(req.Request.URL.RawQuery) > 0 {
		url += "?" + req.Request.URL.RawQuery
	}

	if len(req.Request.URL.Fragment) > 0 {
		url += "#" + req.Request.URL.Fragment
	}

	http.Redirect(
		req.Writer,
		req.Request,
		url,
		http.StatusPermanentRedirect,
	)
}

// MassaTLDInterceptor intercepts request for web on-chain.
func MassaTLDInterceptor(req *interceptor.Interceptor) *interceptor.Interceptor {
	if req == nil {
		return nil
	}

	massaIndex := strings.Index(req.Request.Host, ".massa")

	if massaIndex > 0 && !strings.HasPrefix(req.Request.Host, "my.massa") {
		// If the request is not https, redirect to https.
		if req.Request.TLS == nil {
			redirectTLDToHTTPS(req)
		}

		handleMassaDomainRequest(req.Writer, req.Request, massaIndex)

		return nil
	}

	return req
}

func RedirectToDefaultResourceInterceptor(req *interceptor.Interceptor) *interceptor.Interceptor {
	if req == nil {
		return nil
	}

	prefixes := []string{"/browse/", "/thyra/"}

	for _, prefix := range prefixes {
		if !strings.HasPrefix(req.Request.URL.Path, prefix) {
			continue
		}

		splited := removeEmptyStrings(strings.Split(req.Request.URL.Path[len(prefix):], "/"))

		if len(splited) == 1 {
			protocol := "https"
			if req.Request.TLS == nil {
				protocol = "http"
			}

			http.Redirect(
				req.Writer,
				req.Request,
				protocol+"://"+req.Request.Host+prefix+splited[0]+"/index.html",
				http.StatusSeeOther,
			)

			return nil
		}
	}

	return req
}
