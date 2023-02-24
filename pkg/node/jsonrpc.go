package node

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"

	"github.com/ybbus/jsonrpc/v3"
)

type withLoggingRoundTripper struct {
	isEnabled        bool
	showResponseBody bool
}

func (t *withLoggingRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	if t.isEnabled {
		requestDump, err := httputil.DumpRequestOut(request, true)
		if err != nil {
			panic(fmt.Errorf("unexpecting that httputil.DumpRequestOut would panic: %w", err))
		}
		//nolint:forbidigo
		fmt.Println("JSON-RPC request\n" + string(requestDump))
	}

	response, err := http.DefaultTransport.RoundTrip(request)
	if err != nil {
		return nil, fmt.Errorf("sending http request '%+v': %w", request, err)
	}

	if t.isEnabled {
		responseDump, err := httputil.DumpResponse(response, t.showResponseBody)
		if err != nil {
			panic(fmt.Errorf("unexpecting that httputil.DumpRequestOut would panic: %w", err))
		}
		//nolint:forbidigo
		fmt.Println("JSON-RPC response\n" + string(responseDump))
	}

	return response, nil
}

type Client struct {
	RPCClient jsonrpc.RPCClient
}

func NewDefaultClient() *Client {
	defaultServer := os.Getenv("MASSA_NODE_URL")

	return NewClient(defaultServer)
}

//nolint:exhaustruct
func NewClient(url string) *Client {
	return &Client{
		RPCClient: jsonrpc.NewClientWithOpts(
			url,
			&jsonrpc.RPCClientOpts{
				HTTPClient: &http.Client{Transport: &withLoggingRoundTripper{
					isEnabled:        os.Getenv("DEBUG_RPC") == "true",
					showResponseBody: os.Getenv("DEBUG_RPC") == "true",
				}},
			}),
	}
}
