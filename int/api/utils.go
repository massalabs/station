package api

import (
	"net/http"
	"path/filepath"

	"github.com/go-openapi/runtime"
)

func contentType(rsc string) map[string]string {
	var contentType map[string]string

	switch filepath.Ext(rsc) {
	case ".css":
		contentType = map[string]string{"Content-Type": "text/css"}
	case ".js":
		contentType = map[string]string{"Content-Type": "text/javascript"}
	case ".html":
		contentType = map[string]string{"Content-Type": "text/html"}
	case ".webp":
		contentType = map[string]string{"Content-Type": "text/webp"}
	case ".png":
		contentType = map[string]string{"Content-Type": "image/png"}
	default:
		contentType = map[string]string{}
	}

	return contentType
}

type CustomResponder struct {
	Body       []byte
	Header     map[string]string
	StatusCode int
}

func (c *CustomResponder) WriteResponse(writer http.ResponseWriter, producer runtime.Producer) {
	writer.WriteHeader(c.StatusCode)

	for k, v := range c.Header {
		writer.Header().Add(k, v)
	}

	_, err := writer.Write(c.Body)
	if err != nil {
		panic(err)
	}
}

func NewCustomResponder(body []byte, header map[string]string, statusCode int) *CustomResponder {
	return &CustomResponder{Body: body, Header: header, StatusCode: statusCode}
}

func NewNotFoundResponder() *CustomResponder {
	return NewCustomResponder(
		[]byte("Page not found"),
		map[string]string{"Content-Type": "text/html"},
		http.StatusNotFound)
}

func NewInternalServerErrorResponder(err error) *CustomResponder {
	return NewCustomResponder(
		[]byte(err.Error()),
		map[string]string{"Content-Type": "text/html"},
		http.StatusInternalServerError)
}
