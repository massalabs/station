package massa

import (
	"testing"

	"github.com/massalabs/station/api/swagger/server/restapi/operations"
	"github.com/massalabs/station/int/config"
	"github.com/stretchr/testify/assert"
)

func TestNodeStatus(t *testing.T) {
	t.Run("should return node status", func(t *testing.T) {
		config := &config.NetworkInfos{
			Network:    "testnet",
			NodeURL:    "http://localhost:8080",
			DNSAddress: "AU123",
			ChainID:    1,
		}
		getNodeHandler := &getNodeHandler{config: config}
		response := getNodeHandler.Handle(operations.NewGetNodeParams())
		responseTypes, ok := response.(*operations.GetNodeOK)
		assert.True(t, ok, "responsePayload is not of type *operations.GetNodeOK")
		responsePayload := responseTypes.Payload
		assert.Equal(t, config.Network, responsePayload.Network)
		assert.Equal(t, config.DNSAddress, responsePayload.DNS)
		assert.Equal(t, config.ChainID, uint64(responsePayload.ChainID))
		assert.Equal(t, config.NodeURL, *responsePayload.URL)
	})
}
