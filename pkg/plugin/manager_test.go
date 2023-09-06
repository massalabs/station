package plugin

import (
	"os"
	"path/filepath"
	"testing"
)

//nolint:funlen
func Test_removePlugin(t *testing.T) {
	testCases := []struct {
		name        string
		plugin      *Plugin
		createFiles []string
		remainFiles []string
	}{
		{
			name: "massa-labs/massa-wallet",
			plugin: &Plugin{
				info: &Information{
					Author: "massa-labs",
					Name:   "massa-wallet",
				},
			},
			createFiles: []string{"wallet_test.yaml", "other_file.txt"},
			remainFiles: []string{"wallet_test.yaml"},
		},
		{
			name: "massa-labs/not-massa-wallet",
			plugin: &Plugin{
				info: &Information{
					Author: "massa-labs",
					Name:   "massa-wallet",
				},
			},
			createFiles: []string{"wallet_test.yaml", "other_file.txt"},
			remainFiles: []string{},
		},
		{
			name: "not-massa-labs/massa-wallet",
			plugin: &Plugin{
				info: &Information{
					Author: "massa-labs",
					Name:   "massa-wallet",
				},
			},
			createFiles: []string{"wallet_test.yaml", "other_file.txt"},
			remainFiles: []string{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// Create a temporary directory for testing.
			tempDir, err := os.MkdirTemp("", "plugin_test")
			if err != nil {
				t.Fatalf("Failed to create temp directory: %v", err)
			}

			testCase.plugin.Path = tempDir

			// Create dummy files.
			for _, file := range testCase.createFiles {
				dummyFile := filepath.Join(tempDir, file)
				if err := os.WriteFile(dummyFile, []byte("test"), 0o600); err != nil {
					t.Fatalf("Failed to create dummy file: %v", err)
				}
			}

			// Call the function to test.
			if err := removePlugin(testCase.plugin); err != nil {
				t.Fatalf("Failed to remove plugin: %v", err)
			}

			// Check if the remaining files still exist.
			for _, file := range testCase.remainFiles {
				dummyFile := filepath.Join(tempDir, file)
				if _, err := os.Stat(dummyFile); os.IsNotExist(err) {
					t.Fatalf("File was deleted: %v", err)
				}
			}

			// Clean up the temporary directory.
			os.RemoveAll(tempDir)
		})
	}
}
