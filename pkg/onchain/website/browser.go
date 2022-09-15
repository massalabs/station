package website

import (
	"net/http"
	"path/filepath"

	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/onchain/storage"
)

// TODO use a local cache to reduce network bandwidth.
func Fetch(c *node.Client, addr string, filename string) ([]byte, error) {
	m, err := storage.Get(c, addr, "massa_web")
	if err != nil {
		return nil, err
	}

	return m[filename], nil
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

func setContentType(rsc string, writer http.ResponseWriter) {
	switch filepath.Ext(rsc) {
	case ".css":
		writer.Header().Set("Content-Type", "text/css")
	case ".js":
		writer.Header().Set("Content-Type", "application/json")
	case ".html":
		writer.Header().Set("Content-Type", "text/html")
	case ".webp":
		writer.Header().Set("Content-Type", "text/webp")
	case ".png":
		writer.Header().Set("Content-Type", "image/png")

	}

}

func Request(writer http.ResponseWriter, reader *http.Request, client *node.Client, address string, resource string) {
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
