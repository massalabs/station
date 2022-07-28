package website

import (
	"errors"
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/massalabs/thyra/pkg/front"
	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/onchain/storage"
	"github.com/massalabs/thyra/pkg/wallet"
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

func handleInitialRequest(w http.ResponseWriter, r *http.Request) {
	addr := r.URL.Query()["url"][0]
	rpcClient := node.NewClient("http://145.239.66.206:33035")
	cookie := &http.Cookie{
		Name:   "ocw",
		Value:  addr,
		MaxAge: 10,
	}
	http.SetCookie(w, cookie)

	body, err := Fetch(rpcClient, addr, "index.html")
	if err != nil {
		panic(err)
	}

	w.Write(body)
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

	_, err = w.Write(body)
	if err != nil {
		panic(err)
	}
}

func handleSubsequentRequest(w http.ResponseWriter, r *http.Request) {
	addr, err := r.Cookie("ocw")
	if err != nil {
		fmt.Println("Error reading cookie")
		panic(err)
	}

	rpcClient := node.NewClient("http://145.239.66.206:33035")

	body, err := Fetch(rpcClient, addr.Value, path.Base(r.URL.Path))
	if err != nil {
		panic(err)
	}

	w.Write(body)
}

func HandlerFunc(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/mgmt") {
			handler.ServeHTTP(w, r)
		} else if strings.Index(r.Host, "webuploader.mythyra.massa") != -1 {
			HandleWebsiteUploaderManagementRequest(w, r)
		} else if strings.Index(r.Host, "wallet.mythyra.massa") != -1 {
			wallet.HandleWalletManagementRequest(w, r)
		} else if strings.Index(r.Host, ".massa") > 0 {
			handleMassaDomainRequest(w, r)
		} else if strings.HasPrefix(r.URL.Path, "/website") {
			handleInitialRequest(w, r)
		}
	})
}

//TODO Manage panic(err)
func HandleWebsiteUploaderManagementRequest(w http.ResponseWriter, r *http.Request) {

	target := r.URL.Path[1:]
	var fileText string
	if strings.Index(target, ".css") > 0 {
		fileText = front.WebsiteCss
		w.Header().Set("Content-Type", "text/css")
	} else if strings.Index(target, ".js") > 0 {
		fileText = front.WebsiteJs
		w.Header().Set("Content-Type", "application/json")
	} else if strings.Index(target, ".html") > 0 {
		fileText = front.WebsiteHtml
		w.Header().Set("Content-Type", "text/html")
	} else if strings.Index(target, ".webp") > 0 {
		fileText = front.Logo_massaWebp
		w.Header().Set("Content-Type", "image/webp")
	}
	_, err := w.Write([]byte(fileText))
	if err != nil {
		panic(err)
	}
}
