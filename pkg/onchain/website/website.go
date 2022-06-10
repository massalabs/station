package website

import (
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/massalabs/thyra/pkg/onchain/storage"
)

func Fetch(addr string, filename string) ([]byte, error) {
	//TODO use a local cache to reduce network bandwidth

	m, err := storage.Get(addr, "2qbtmxh5pD3TH3McFmZWxvKLTyz2SKDYFSRL8ngQBJ4R6f3Duw")
	if err != nil {
		return nil, err
	}

	msg := ""
	for k := range m {
		msg = msg + k + ", "
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

func HandlerFunc(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/website") {
			handleInitialRequest(w, r)
		} else if strings.Contains(path.Base(r.URL.Path), ".") {
			handleSubsequentRequest(w, r)
		} else {
			handler.ServeHTTP(w, r)
		}
	})

}
