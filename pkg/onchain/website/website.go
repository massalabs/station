package website

import (
	"errors"
	"fmt"
	"net/http"
	"path"
	"strings"

	"lukechampine.com/blake3"

	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/node/base58"
	"github.com/massalabs/thyra/pkg/node/ledger"
	"github.com/massalabs/thyra/pkg/onchain/storage"
)

func Resolve(name string) (string, error) {
	digest := blake3.Sum256([]byte("record" + name))
	key := base58.CheckEncode(digest[:])

	address := "A12o8tw83tDrA52Lju9BUDDodAhtUp4scHtYr8Fj4iwhDTuWZqHZ"

	c := node.NewClient("https://test.massa.net/api/v2")

	content, err := ledger.Addresses(c, []string{address})
	if err != nil {
		return "", err
	}

	val, ok := content[0].Info.Datastore[key]
	if ok {
		return string(val), nil
	}

	return "", errors.New("name not found")
}

func Fetch(addr string, filename string) ([]byte, error) {
	// TODO use a local cache to reduce network bandwidth

	m, err := storage.Get(addr, "2qbtmxh5pD3TH3McFmZWxvKLTyz2SKDYFSRL8ngQBJ4R6f3Duw")
	if err != nil {
		return nil, err
	}

	return m[filename], nil
}

func handleInitialRequest(w http.ResponseWriter, r *http.Request) {
	addr := r.URL.Query()["url"][0]

	cookie := &http.Cookie{
		Name:   "ocw",
		Value:  addr,
		MaxAge: 10,
	}
	http.SetCookie(w, cookie)

	body, err := Fetch(addr, "index.html")
	if err != nil {
		panic(err)
	}

	w.Write(body)
}

func handleSubsequentRequest(w http.ResponseWriter, r *http.Request) {
	addr, err := r.Cookie("ocw")
	if err != nil {
		fmt.Println("Error reading cookie")
		panic(err)
	}

	body, err := Fetch(addr.Value, path.Base(r.URL.Path))
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

	addr, err := Resolve(name)
	if err != nil {
		panic(err)
	}

	fmt.Println("Name resolved: " + name + ".massa => " + addr)

	var target string
	if r.URL.Path == "/" {
		target = "index.html"
	} else {
		target = r.URL.Path[1:]
	}

	body, err := Fetch(addr, target)
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
	// fmt.Println(target, body)

	w.Write(body)
}

func HandlerFunc(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Index(r.Host, ".massa") > 0 {
			handleMassaDomainRequest(w, r)
		} else if strings.HasPrefix(r.URL.Path, "/website") {
			handleInitialRequest(w, r)
		} else if strings.Contains(path.Base(r.URL.Path), ".") {
			handleSubsequentRequest(w, r)
		} else {
			handler.ServeHTTP(w, r)
		}
	})
}
