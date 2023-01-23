package plugin

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/massalabs/thyra/api/interceptor"
)

func NewAPIPHandler(manager *Manager) *APIHandler {
	return &APIHandler{manager: manager}
}

type APIHandler struct {
	manager *Manager
}

func (h *APIHandler) Handle(writer http.ResponseWriter, reader *http.Request, pluginAuthor string, pluginName string) {
	alias := fmt.Sprintf("%s/%s", pluginAuthor, pluginName)

	plugin := h.manager.PluginByAlias(alias)

	if plugin == nil {
		writer.WriteHeader(http.StatusNotFound)

		return
	}

	plugin.ReverseProxy().ServeHTTP(writer, reader)
}

//nolint:gochecknoglobals
var Handler APIHandler

const endpointPattern = "/thyra/plugin/"

// Interceptor intercepts requests for plugins.
// The endpoint is expected to have the following structure:
// /thyra/plugin/{author-name}/{plugin-name}/{plugin-endpoint}...
func Interceptor(req *interceptor.Interceptor) *interceptor.Interceptor {
	if req == nil {
		return nil
	}

	isMyMassa := strings.HasPrefix(req.Request.Host, "my.massa")
	indexPluginEndpoint := strings.Index(req.Request.RequestURI, endpointPattern)

	if isMyMassa && indexPluginEndpoint > -1 {
		indexAuthorName := indexPluginEndpoint + len(endpointPattern)
		indexPluginName := indexAuthorName + strings.Index(req.Request.RequestURI[indexAuthorName:], "/") + 1
		indexPluginEndpoint := indexPluginName + strings.Index(req.Request.RequestURI[indexPluginName:], "/") + 1

		authorName, err := url.QueryUnescape(req.Request.RequestURI[indexAuthorName : indexPluginName-1])
		if err != nil {
			log.Fatal(err)

			return nil
		}

		pluginName, err := url.QueryUnescape(req.Request.RequestURI[indexPluginName : indexPluginEndpoint-1])
		if err != nil {
			log.Fatal(err)

			return nil
		}

		Handler.Handle(
			req.Writer, req.Request,
			authorName, pluginName,
		)

		return nil
	}

	return req
}

// modifyRequest rewrite the incoming request URL to match what the pkugin is expecting to receive.
// All the `/thyra/plugin/{author-name}/{plugin-name}` template is removed.
func modifyRequest(req *http.Request) {
	urlExternal := req.URL.String()

	endpointPatternLength := len(endpointPattern)

	// the url has the following format:
	// 		http://127.0.0.1:1234/thyra/plugin/massalabs/hello%20world/web/index.html?name=Massalabs
	// The idea is to rewrite url to remove: /thyra/plugin/massalabs/hello%20world
	// prefixBegin is for the first slash.
	// indexPluginName is for the first char after the slash after massalabs
	// prefixEnd is for slash after hello%20world
	prefixBegin := strings.Index(urlExternal, endpointPattern)

	indexPluginName := prefixBegin + endpointPatternLength +
		strings.Index(urlExternal[prefixBegin+endpointPatternLength:], "/") + 1

	prefixEnd := indexPluginName + strings.Index(urlExternal[indexPluginName:], "/")

	urlRewritten := urlExternal[:prefixBegin] + urlExternal[prefixEnd:]

	newURL, err := url.Parse(urlRewritten)
	if err != nil {
		panic(err)
	}

	req.URL = newURL
}
