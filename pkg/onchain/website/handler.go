package website

import (
	"net/http"
	"strings"

	"github.com/massalabs/thyra/api/interceptor"
	"github.com/massalabs/thyra/pkg/config"
	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/onchain/dns"
)

func handleMassaDomainRequest(writer http.ResponseWriter, reader *http.Request, index int, config config.AppConfig) {
	name := reader.Host[:index]

	rpcClient := node.NewClient(config.NodeURL)

	addr, err := dns.Resolve(config, rpcClient, name)
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

// MassaTLDInterceptor intercepts request for web on-chain.
func MassaTLDInterceptor(req *interceptor.Interceptor, appConfig config.AppConfig) *interceptor.Interceptor {
	if req == nil {
		return nil
	}

	massaIndex := strings.Index(req.Request.Host, ".massa")

	if massaIndex > 0 && !strings.HasPrefix(req.Request.Host, config.MassaStationURL) {
		handleMassaDomainRequest(req.Writer, req.Request, massaIndex, appConfig)

		return nil
	}

	return req
}

func RedirectToDefaultResourceInterceptor(req *interceptor.Interceptor) *interceptor.Interceptor {
	if req == nil {
		return nil
	}

	prefixes := []string{"/home", "/search", "/websiteUploader", "/thyra/plugin-manager"}

	for _, prefix := range prefixes {
		if !strings.HasPrefix(req.Request.URL.Path, prefix) {
			continue
		}

		// before we had something like /thyra/websiteUploader, which needs to be converted
		// to an array [websiteUploader] in this specific case the len is 1, and we proceed
		// adding "/.html" to  "websiteuploader"
		// Now i want to keep same functionality but we don't have anymore /thyra, but instead
		// only /websiteUploader, it will be converted to [] only in this case len is 0,
		// i can start adding  "/.html"
		// it's working but can be refactored
		// TODO : specific PR to refector this interceptor
		splited := removeEmptyStrings(strings.Split(req.Request.URL.Path[len(prefix):], "/"))

		if len(splited) == 0 {
			protocol := "https"
			if req.Request.TLS == nil {
				protocol = "http"
			}

			http.Redirect(
				req.Writer,
				req.Request,
				protocol+"://"+req.Request.Host+prefix+"/index.html",
				http.StatusSeeOther,
			)

			return nil
		}
	}

	return req
}
