package logs

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station/api/swagger/server/models"
	"github.com/massalabs/station/api/swagger/server/restapi/operations"
)

const (
	errorCodeReadLogs = "Logs-0001"

	// timestampLayout is the layout of zapcore.ISO8601TimeEncoder, which starts each log entry.
	timestampLayout = "2006-01-02T15:04:05.000Z0700"
)

// ansiEscapeRegexp matches the terminal color codes written by the logger level encoder.
var ansiEscapeRegexp = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// NewGetLogsHandler returns a handler serving the content of the log file at the given path.
// sessionStart is when the current MassaStation session started logging.
func NewGetLogsHandler(logFilePath string, sessionStart time.Time) operations.GetMassaStationLogsHandler {
	return operations.GetMassaStationLogsHandlerFunc(func(params operations.GetMassaStationLogsParams) middleware.Responder {
		content, err := os.ReadFile(logFilePath)
		if err != nil {
			return operations.NewGetMassaStationLogsInternalServerError().WithPayload(&models.Error{
				Code:    errorCodeReadLogs,
				Message: fmt.Sprintf("read log file: %s", err),
			})
		}

		logs := ansiEscapeRegexp.ReplaceAllString(string(content), "")

		if params.CurrentSession != nil && *params.CurrentSession {
			logs = logsSince(logs, sessionStart)
		}

		if params.Lines != nil {
			logs = lastLines(logs, int(*params.Lines))
		}

		return operations.NewGetMassaStationLogsOK().WithPayload(logs)
	})
}

// lastLines returns the last n lines of the given text.
func lastLines(text string, n int) string {
	lines := strings.SplitAfter(strings.TrimSuffix(text, "\n"), "\n")
	if len(lines) <= n {
		return text
	}

	return strings.Join(lines[len(lines)-n:], "") + "\n"
}

// logsSince returns the logs starting at the first entry written at or after the given time.
// Lines without timestamp (e.g. stack traces) belong to the preceding entry.
func logsSince(text string, start time.Time) string {
	// Log timestamps have a millisecond precision.
	start = start.Truncate(time.Millisecond)

	offset := 0
	for offset < len(text) {
		line, _, _ := strings.Cut(text[offset:], "\n")
		timestamp, _, _ := strings.Cut(line, "\t")

		if t, err := time.Parse(timestampLayout, timestamp); err == nil && !t.Before(start) {
			return text[offset:]
		}

		offset += len(line) + 1
	}

	return ""
}
