package massa

import (
	"testing"

	"github.com/massalabs/station/int/config"
)

func TestNodeStatus(t *testing.T) {
	t.Run("should return node status", func(t *testing.T) {
		// Reset singleton for testing
		config.ResetConfigManager()

		// Note: This test would require actual network setup or mocking
		// For now, we'll skip it as it needs the singleton to be properly initialized
		// with network data, which requires file system operations
		t.Skip("Integration test - requires full config manager setup")
	})
}
