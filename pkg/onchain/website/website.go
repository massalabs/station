package website

import (
	"errors"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/massalabs/thyra/pkg/front"
	fwallet "github.com/massalabs/thyra/pkg/front/wallet"
	"github.com/massalabs/thyra/pkg/front/website"
	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/onchain/storage"
)

func Resolve(client *node.Client, name string) (string, error) {
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

func handleAPIRequest(w http.ResponseWriter, r *http.Request) {
	prefixSize := len("/website/")
	if len(r.URL.Path) < prefixSize {
		pathNotFound(w, r)

		return
	}

	splited := removeEmptyStrings(strings.Split(r.URL.Path[prefixSize:], "/"))

	switch len(splited) {
	// no resource, only an address is present
	case 1:
		http.Redirect(w, r, "http://"+r.Host+"/website/"+splited[0]+"/"+"index.html", http.StatusSeeOther)
	// address and resource are present
	case 2:
		c := node.NewClient("http://145.239.66.206:33035")
		Request(w, r, c, splited[0], splited[1])
	default:
		pathNotFound(w, r)
	}
}

func handleMassaDomainRequest(w http.ResponseWriter, r *http.Request, index int) {
	name := r.Host[:index]

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

	Request(w, r, rpcClient, addr, target)
}

func HandlerFunc(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		massaTLD := strings.Index(r.Host, ".massa")

		if massaTLD > 0 {
			handleMassaDomainRequest(w, r, massaTLD)
		} else if strings.HasPrefix(r.URL.Path, "/website") {
			handleAPIRequest(w, r)
		} else if strings.HasPrefix(r.URL.Path, "/webuploader.mythyra.massa") {
			HandleWebsiteUploaderManagementRequest(w, r)
		} else if strings.HasPrefix(r.URL.Path, "/wallet.mythyra.massa") {
			HandleWalletManagementRequest(w, r)
		} else {
			handler.ServeHTTP(w, r)
		}
	})
}

func HandleWebsiteUploaderManagementRequest(w http.ResponseWriter, r *http.Request) {
	resource := r.URL.Path[1:]

	var fileText string

	switch resource {
	case "website.css":
		fileText = website.CSS
	case "website.js":
		fileText = website.JS
	case "website.html":
		fileText = website.HTML
	case "logo.webp.css":
		fileText = front.Logo
	}

	setContentType(resource, w)

	_, err := w.Write([]byte(fileText))
	if err != nil {
		panic(err)
	}
}

func HandleWalletManagementRequest(w http.ResponseWriter, r *http.Request) {
	resource := r.URL.Path[1:]

	var fileText string

	switch resource {
	case "website.css":
		fileText = fwallet.CSS
	case "website.js":
		fileText = fwallet.JS
	case "website.html":
		fileText = fwallet.HTML
	case "logo.webp.css":
		fileText = front.Logo
	}

	setContentType(resource, w)

	_, err := w.Write([]byte(fileText))
	if err != nil {
		panic(err)
	}
}
