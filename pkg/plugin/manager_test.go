package plugin

import (
	"os"
	"path/filepath"
	"testing"
)

func Test_removePlugin(t *testing.T) {
	testCases := []struct {
		name        string
		plugin      *Plugin
		createFiles []string
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

			// Clean up the temporary directory.
			os.RemoveAll(tempDir)
		})
	}
}
