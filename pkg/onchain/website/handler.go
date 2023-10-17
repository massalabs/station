package website

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/massalabs/station/api/interceptor"
	"github.com/massalabs/station/int/config"
	"github.com/massalabs/station/pkg/dnshelper"
	"github.com/massalabs/station/pkg/logger"
	"github.com/massalabs/station/pkg/node"
	"github.com/massalabs/station/pkg/onchain/dns"
)

func handleMassaDomainRequest(
	writer http.ResponseWriter,
	reader *http.Request,
	index int,
	config config.NetworkInfos,
	configDir string,
) error {
	name := reader.Host[:index]

	rpcClient := node.NewClient(config.NodeURL)

	addr, err := dns.Resolve(config, rpcClient, name)
	if err != nil {
		return fmt.Errorf("resolving '%s': %w", name, err)
	}

	var target string
	if reader.URL.Path == "/" {
		target = "index.html"
	} else {
		target = reader.URL.Path[1:]
	}

	return Request(writer, reader, rpcClient, addr, target, configDir)
}

// MassaTLDInterceptor intercepts request for web on-chain.
func MassaTLDInterceptor(
	req *interceptor.Interceptor,
	appConfig config.NetworkInfos,
	configDir string,
) *interceptor.Interceptor {
	if req == nil {
		return nil
	}

	massaIndex := strings.Index(req.Request.Host, ".massa")

	if massaIndex > 0 && !strings.HasPrefix(req.Request.Host, config.MassaStationURL) {
		// This is mandatory to avoid CORS issues.
		req.Writer.Header().Set("Access-Control-Allow-Origin", "*")

		err := handleMassaDomainRequest(req.Writer, req.Request, massaIndex, appConfig, configDir)
		if err != nil {
			isFaviconError := strings.Contains(err.Error(), dnshelper.FaviconIcon)

			// avoid spamming logs with favicon.ico fetching errors
			if isFaviconError {
				logger.Warn("handling massa domain request:", err)
			}

			http.NotFound(req.Writer, req.Request)

			return nil
		}

		return nil
	}

	return req
}
