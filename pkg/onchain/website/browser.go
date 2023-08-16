package website

import (
	"fmt"
	"net/http"

	"github.com/massalabs/station/int/api/utils"
	"github.com/massalabs/station/pkg/node"
	"github.com/massalabs/station/pkg/onchain/storage"
)

func Fetch(c *node.Client, addr string, filename string) ([]byte, error) {
	m, err := storage.Get(c, addr)
	if err != nil {
		return nil, fmt.Errorf("fetching the '%s' web resource at '%s': %w", filename, addr, err)
	}

	return m[filename], nil
}

func setContentType(file string, writer http.ResponseWriter) {
	writer.Header().Set(utils.ContentTypeHeader, utils.ContentType(file)[utils.ContentTypeHeader])
}

func Request(writer http.ResponseWriter, _ *http.Request, client *node.Client, address string, resource string) error {
	body, err := Fetch(client, address, resource)
	if err != nil {
		return fmt.Errorf("fetching the '%s' web resource at '%s': %w", resource, address, err)
	}

	setContentType(resource, writer)

	_, err = writer.Write(body)
	if err != nil {
		return fmt.Errorf("writing the '%s' web resource at '%s': %w", resource, address, err)
	}

	return nil
}
