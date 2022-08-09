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
	var r []string

	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}

	return r
}

func pathNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "application/json")

	_, err := w.Write([]byte("{\"code\":404,\"message\":\"path " + r.URL.Path + " was not found\"}"))
	if err != nil {
		panic(err)
	}
}

func setContentType(rsc string, w http.ResponseWriter) {
	switch filepath.Ext(rsc) {
	case ".css":
		w.Header().Set("Content-Type", "text/css")
	case ".js":
		w.Header().Set("Content-Type", "application/json")
	case ".html":
		w.Header().Set("Content-Type", "text/html")
	case ".webp":
		w.Header().Set("Content-Type", "text/webp")
	case ".png":
		w.Header().Set("Content-Type", "image/png")
	}
}

func Request(w http.ResponseWriter, r *http.Request, c *node.Client, address string, resource string) {
	body, err := Fetch(c, address, resource)
	if err != nil {
		panic(err)
	}
	setContentType(resource, w)

	_, err = w.Write(body)
	if err != nil {
		panic(err)
	}
}
