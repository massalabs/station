package website

import (
	"fmt"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/onchain/storage"
)

func Fetch(c *node.Client, addr string, filename string) ([]byte, error) {
	m, err := storage.Get(c, addr)
	if err != nil {
		return nil, fmt.Errorf("fetching the '%s' web resource at '%s': %w", filename, addr, err)
	}

	return m[filename], nil
}

func setContentType(file string, writer http.ResponseWriter) {
	ctype := mime.TypeByExtension(filepath.Ext(file))

	if ctype == "" {
		ctype = "text/plain"
	}

	writer.Header().Set("Content-Type", ctype)
}

func Request(writer http.ResponseWriter, _ *http.Request, client *node.Client, address string, resource string) {
	body, err := Fetch(client, address, resource)
	if err != nil {
		panic(err)
	}

	setContentType(resource, writer)

	_, err = writer.Write(body)
	if err != nil {
		panic(err)
	}
}
