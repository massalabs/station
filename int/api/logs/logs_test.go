package logs

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-openapi/runtime"
	"github.com/massalabs/station/api/swagger/server/restapi/operations"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const logContent = "2026-10-29T10:00:00.000+0100\t\x1b[34mINFO\x1b[0m\tfirst\n" +
	"2026-10-29T10:00:01.000+0100\t\x1b[33mWARN\x1b[0m\tsecond\n" +
	"2026-10-29T10:00:02.000+0100\t\x1b[31mERROR\x1b[0m\tthird\n" +
	"main.main\n" +
	"\t/station/main.go:42\n"

const thirdEntry = "2026-10-29T10:00:02.000+0100\tERROR\tthird\n" +
	"main.main\n" +
	"\t/station/main.go:42\n"

func TestGetLogsHandler(t *testing.T) {
	logFilePath := filepath.Join(t.TempDir(), "massastation.log")
	require.NoError(t, os.WriteFile(logFilePath, []byte(logContent), 0o600))

	int64Ptr := func(v int64) *int64 { return &v }
	boolPtr := func(v bool) *bool { return &v }

	// Between the second and third entries, with a sub-millisecond part ignored by log timestamps.
	sessionStart := time.Date(2026, 10, 29, 9, 0, 1, 500_000_000, time.UTC)

	tests := []struct {
		name           string
		lines          *int64
		currentSession *bool
		sessionStart   time.Time
		want           string
	}{
		{
			name: "whole file without color codes",
			want: "2026-10-29T10:00:00.000+0100\tINFO\tfirst\n" +
				"2026-10-29T10:00:01.000+0100\tWARN\tsecond\n" +
				thirdEntry,
		},
		{
			name:  "last lines",
			lines: int64Ptr(2),
			want:  "main.main\n\t/station/main.go:42\n",
		},
		{
			name:  "more lines than available",
			lines: int64Ptr(10),
			want: "2026-10-29T10:00:00.000+0100\tINFO\tfirst\n" +
				"2026-10-29T10:00:01.000+0100\tWARN\tsecond\n" +
				thirdEntry,
		},
		{
			name:           "current session",
			currentSession: boolPtr(true),
			sessionStart:   sessionStart,
			want:           thirdEntry,
		},
		{
			name:           "current session starting at the same millisecond as an entry",
			currentSession: boolPtr(true),
			sessionStart:   time.Date(2026, 10, 29, 9, 0, 2, 400_000, time.UTC),
			want:           thirdEntry,
		},
		{
			name:           "current session without logs yet",
			currentSession: boolPtr(true),
			sessionStart:   time.Date(2026, 10, 29, 9, 0, 3, 0, time.UTC),
			want:           "",
		},
		{
			name:           "current session disabled",
			currentSession: boolPtr(false),
			sessionStart:   sessionStart,
			want: "2026-10-29T10:00:00.000+0100\tINFO\tfirst\n" +
				"2026-10-29T10:00:01.000+0100\tWARN\tsecond\n" +
				thirdEntry,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewGetLogsHandler(logFilePath, tt.sessionStart)
			resp := handler.Handle(operations.GetMassaStationLogsParams{Lines: tt.lines, CurrentSession: tt.currentSession})

			okResp, ok := resp.(*operations.GetMassaStationLogsOK)
			require.True(t, ok, "unexpected response type %T", resp)
			assert.Equal(t, tt.want, okResp.Payload)
		})
	}
}

func TestGetLogsHandlerMissingFile(t *testing.T) {
	handler := NewGetLogsHandler(filepath.Join(t.TempDir(), "missing.log"), time.Now())
	resp := handler.Handle(operations.GetMassaStationLogsParams{})

	_, ok := resp.(*operations.GetMassaStationLogsInternalServerError)
	require.True(t, ok, "unexpected response type %T", resp)

	recorder := httptest.NewRecorder()
	resp.WriteResponse(recorder, runtime.TextProducer())
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	assert.Contains(t, recorder.Body.String(), errorCodeReadLogs)
}
