package website

import (
	"errors"
	"net/http"
	"strings"

	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/onchain/storage"
)

func Resolve(client *node.Client, name string) (string, error) {
	//digest := blake3.Sum256([]byte("record" + name))
	//key := base58.CheckEncode(digest[:])

	address := "A12ew8eiCS7wnY8SkUdwBgDkdD5qwmbJgkJvYLCvVjWWdoFJJLvW"

	const dnsPrefix = "record"

	entry, err := node.DatastoreEntry(client, address, dnsPrefix+name)
	if err != nil {
		return "", err
	}

	if len(entry.CandidateValue) == 0 {
		return "", errors.New("name not found")
	}

	return string(entry.CandidateValue), nil
}

// TODO use a local cache to reduce network bandwidth.
func Fetch(c *node.Client, addr string, filename string) ([]byte, error) {
	m, err := storage.Get(c, addr, "massa_web")
	if err != nil {
		return nil, err
	}

	return m[filename], nil
}

func handleMassaDomainRequest(w http.ResponseWriter, r *http.Request) {
	i := strings.Index(r.Host, ".massa")
	if i < 0 {
		panic("no .massa in URL")
	}

	name := r.Host[:i]

	rpcClient := node.NewClient("http://145.239.66.206:33035")

	addr, err := Resolve(rpcClient, name)
	if err != nil {
		panic(err)
	}

	var target string
	if r.URL.Path == "/" {
		target = "index.html"
	} else {
		target = r.URL.Path[1:]
	}

	body, err := Fetch(rpcClient, addr, target)
	if err != nil {
		panic(err)
	}

	if strings.Index(target, ".css") > 0 {
		w.Header().Set("Content-Type", "text/css")
	} else if strings.Index(target, ".js") > 0 {
		w.Header().Set("Content-Type", "application/json")
	} else if strings.Index(target, ".html") > 0 {
		w.Header().Set("Content-Type", "text/html")
	}

	w.Write(body)
}

func HandlerFunc(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Index(r.Host, ".massa") > 0 {
			handleMassaDomainRequest(w, r)
		} else {
			handler.ServeHTTP(w, r)
		}
	})
}
