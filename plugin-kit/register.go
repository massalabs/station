package register

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
)

const (
	MassaStationBaseURL          = "http://station.massa"
	PluginManagerEndpoint        = "plugin-manager/register"
	MassaStationRegisterEndpoint = MassaStationBaseURL + "/" + PluginManagerEndpoint
	StandaloneEnvVar             = "STANDALONE"
)

// ErrMissingPluginID is returned by RegisterPlugin when the plugin is not started with its ID as first argument.
var ErrMissingPluginID = errors.New("program must be started with a correlationID command line argument")

type registerBody struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

func RegisterPlugin(listener net.Listener) error {
	minimumNumberOfCLIArgument := 2

	if len(os.Args) < minimumNumberOfCLIArgument {
		return ErrMissingPluginID
	}

	return register(MassaStationRegisterEndpoint, os.Args[1], listener.Addr())
}

func register(
	endpoint string,
	pluginID string,
	socket net.Addr,
) error {
	reg := registerBody{
		ID:  pluginID,
		URL: "http://" + socket.String(),
	}

	body, err := json.Marshal(reg)
	if err != nil {
		return fmt.Errorf("while marshaling register body: %w", err)
	}

	request, err := http.NewRequestWithContext(context.Background(), http.MethodPost, endpoint, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("while creating register request: %w", err)
	}

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	request.Header.Set("Origin", "localhost")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("while doing register request: %w", err)
	}

	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(response.Body)

		var data interface{}

		err = json.Unmarshal(body, &data)
		if err != nil {
			data = fmt.Errorf("unable to decode json response: %w", err)
		}

		return fmt.Errorf("plugin registry failed: HTTP status: %d, HTTP body: %v", response.StatusCode, data)
	}

	return nil
}
