package logger

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCloseWithPipedStdout(t *testing.T) {
	reader, writer, err := os.Pipe()
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = reader.Close()
		_ = writer.Close()
	})

	stdout := os.Stdout
	os.Stdout = writer

	t.Cleanup(func() { os.Stdout = stdout })

	logger, err := New(filepath.Join(t.TempDir(), "test.log"))
	require.NoError(t, err)

	logger.Info("message")

	require.NoError(t, logger.Close())
}
