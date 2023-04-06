package plugin

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/massalabs/thyra/api/interceptor"
)

func NewAPIHandler(manager *Manager) *APIHandler {
	return &APIHandler{manager: manager}
}

type APIHandler struct {
	manager *Manager
}

func (h *APIHandler) Handle(writer http.ResponseWriter, reader *http.Request, pluginAuthor string, pluginName string) {
	alias := Alias(pluginAuthor, pluginName)

	plugin, err := h.manager.PluginByAlias(alias)
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprint(writer, err)

		return
	}

	plugin.ReverseProxy().ServeHTTP(writer, reader)
}

//nolint:gochecknoglobals
var Handler APIHandler

const EndpointPattern = "/thyra/plugin/"

type endpointContent struct {
	pluginAuthor string
	pluginName   string
	subURI       string
}

func splitEndpoint(uri string) *endpointContent {
	// ["", "thyra", "plugin", "{author-name}", "{plugin-name}", ...]
	exploded := strings.Split(uri, "/")

	return &endpointContent{
		pluginAuthor: FormatTextForURL(exploded[3]),
		pluginName:   FormatTextForURL(exploded[4]),
		subURI:       "/" + strings.Join(exploded[5:], "/"),
	}
}

// Interceptor intercepts requests for plugins.
// The endpoint is expected to have the following structure:
// /thyra/plugin/{author-name}/{plugin-name}/{plugin-endpoint}...
func Interceptor(req *interceptor.Interceptor) *interceptor.Interceptor {
	if req == nil {
		return nil
	}

	isMyMassa := strings.HasPrefix(req.Request.Host, "my.massa")
	indexPluginEndpoint := strings.Index(req.Request.RequestURI, EndpointPattern)

	if isMyMassa && indexPluginEndpoint > -1 {
		endpoint := splitEndpoint(req.Request.RequestURI)

		authorName, err := url.QueryUnescape(endpoint.pluginAuthor)
		if err != nil {
			log.Fatal(err)

			return nil
		}

		pluginName, err := url.QueryUnescape(endpoint.pluginName)
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

// modifyRequest rewrite the incoming request URL to match what the plugin is expecting to receive.
// All the `/thyra/plugin/{author-name}/{plugin-name}` template is removed.
func modifyRequest(req *http.Request) {
	urlExternal := req.URL.String()

	// the url has the following format:
	// 		http://127.0.0.1:1234/thyra/plugin/massalabs/hello-world/web/index.html?name=Massalabs
	// The idea is to rewrite url to remove: /thyra/plugin/massalabs/hello-world

	index := strings.Index(urlExternal, EndpointPattern)

	endpoint := splitEndpoint(urlExternal[index:])

	urlRewritten := urlExternal[:index] + endpoint.subURI

	urlNew, err := url.Parse(urlRewritten)
	if err != nil {
		panic(err)
	}

	req.URL = urlNew
}
