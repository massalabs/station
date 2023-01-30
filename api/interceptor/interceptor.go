package interceptor

import "net/http"

// Interceptor is an HTTP request interceptor.
// An interceptor works the following way:
// if it recognizes the request, it treats it and return nil,
// if it doesn't recognize the request, it returns the given interceptor.
type Interceptor struct {
	Writer  http.ResponseWriter
	Request *http.Request
}
