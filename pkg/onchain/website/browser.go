package website

import (
	"fmt"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/onchain/storage"
)

func Fetch(c *node.Client, addr string, filename string) ([]byte, bool, int, error) {
	m, err := storage.Get(c, addr)
	if err != nil {
		return nil, false, 500, fmt.Errorf("fetching the '%s' web resource at '%s': %w", filename, addr, err)
	}

	if _, ok := m[filename]; !ok {
		// if the file is not found, try to find 404.html
		if _, ok := m["404.html"]; ok {
			return m["404.html"], false, 404, nil
		}

		// otherwise, let the client handle the route
		// return true to indicate that the route is dynamic
		return m["index.html"], true, 200, nil
	}

	return m[filename], false, 200, nil
}

func removeEmptyStrings(s []string) []string {
	var result []string

	for _, str := range s {
		if str != "" {
			result = append(result, str)
		}
	}

	return result
}

func setContentType(file string, writer http.ResponseWriter) {
	ctype := mime.TypeByExtension(filepath.Ext(file))

	if ctype == "" {
		ctype = "text/plain"
	}

	writer.Header().Set("Content-Type", ctype)
}

func Request(writer http.ResponseWriter, reader *http.Request, client *node.Client, address string, resource string) {
	body, dynamic, status, err := Fetch(client, address, resource)
	if err != nil {
		panic(err)
	}

	if dynamic {
		writer.Header().Set("Content-Type", "text/html")
	} else {
		setContentType(resource, writer)
	}

	writer.WriteHeader(status)

	// add Access-Control-Allow-Origin header (for CORS)
	writer.Header().Set("Access-Control-Allow-Origin", "*")

	_, err = writer.Write(body)
	if err != nil {
		panic(err)
	}
}
