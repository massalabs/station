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

	/*msg := ""
	for k := range m {
		msg = msg + k + ", "
	}

	fmt.Println("files", msg)*/

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

var dns = map[string]string{
	"flappy": "A1aMywGBgBywiL6WcbKR4ugxoBtdP9P3waBVi5e713uvj7F1DJL",
	"blog":   "A1NjcatuB6SLecX8xQp8nF3xMr6eW3YDQ2a9i7gaetK9SskJrBi",
}

func handleMassaDomainRequest(w http.ResponseWriter, r *http.Request) {
	i := strings.Index(r.Host, ".massa")
	if i < 0 {
		panic("no .massa in URL")
	}

	name := r.Host[:i]

	addr, ok := dns[name]
	if !ok {
		panic("following name not resolved " + name)
	}

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
	//fmt.Println(target, body)

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
