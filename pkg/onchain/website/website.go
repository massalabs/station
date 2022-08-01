package website

import (
	"errors"
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/massalabs/thyra/pkg/front"
	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/node/getters"
	"github.com/massalabs/thyra/pkg/onchain/storage"
	"github.com/massalabs/thyra/pkg/wallet"
)

func Resolve(client *node.Client, name string) (string, error) {
	//digest := blake3.Sum256([]byte("record" + name))
	//key := base58.CheckEncode(digest[:])

	address := "A12ew8eiCS7wnY8SkUdwBgDkdD5qwmbJgkJvYLCvVjWWdoFJJLvW"

	const dnsPrefix = "record"

	entry, err := getters.DatastoreEntry(client, address, dnsPrefix+name)
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
	rpcClient := node.NewClient()
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

//TO DO TO Remove
func handleMassaDomainRequest(w http.ResponseWriter, r *http.Request) {
	// i := strings.Index(r.Host, ".massa")
	// if i < 0 {
	// 	panic("no .massa in URL")
	// }

	// name := r.Host[:i]

	addr := r.URL.Query().Get("url")
	rpcClient := node.NewClient()

	// addr, err := Resolve(rpcClient, name)
	// if err != nil {
	// 	panic(err)
	// }

	var target string
	if r.URL.Path == "/" {
		target = "index.html"
	} else {
		target = strings.Split(r.URL.Path, "/")[2]
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

	rpcClient := node.NewClient()

	body, err := Fetch(rpcClient, addr.Value, path.Base(r.URL.Path))
	if err != nil {
		panic(err)
	}

	w.Write(body)
}

//TO REWORK
func HandlerFunc(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/mgmt/wallet") {
			handler.ServeHTTP(w, r)
		} else if strings.HasPrefix(r.URL.Path, "/website") {
			handleMassaDomainRequest(w, r)
		} else if strings.HasPrefix(r.URL.Path, "/uploadWeb") && r.Method == "POST" {
			CreateWebsiteDeployer(w, r)
		} else if strings.HasPrefix(r.URL.Path, "/uploadWeb") && r.Method == "GET" {
			RefreshDeployers(w, r)
		} else if strings.HasPrefix(r.URL.Path, "/uploadWeb") && r.Method == "PUT" {
			UploadWebsite(w, r)
		} else if strings.HasPrefix(r.Host, "webuploader.mythyra.massa") {
			HandleWebsiteUploaderManagementRequest(w, r)
		} else if strings.HasPrefix(r.Host, "wallet.mythyra.massa") {
			wallet.HandleWalletManagementRequest(w, r)
		} else {
			handler.ServeHTTP(w, r)
		}
		// else if strings.Index(r.Host, ".massa") > 0 {
		// 	fmt.Println("aca")
		// 	handleMassaDomainRequest(w, r)
		// }

	})
}

//TODO Manage panic(err)
func HandleWebsiteUploaderManagementRequest(w http.ResponseWriter, r *http.Request) {

	target := r.URL.Path[1:]
	var fileText string
	if strings.Index(target, ".css") > 0 {
		fileText = front.PageCss
		w.Header().Set("Content-Type", "text/css")
	} else if strings.Index(target, ".js") > 0 {
		fileText = front.PageJs
		w.Header().Set("Content-Type", "application/json")
	} else if strings.Index(target, ".html") > 0 {
		fileText = front.PageHtml
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
