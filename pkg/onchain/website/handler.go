package website

import (
	"net/http"
	"strings"

	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/onchain/dns"
)

func handleMassaDomainRequest(writer http.ResponseWriter, reader *http.Request, index int) {
	name := reader.Host[:index]

	rpcClient := node.NewDefaultClient()

	addr, err := dns.Resolve(rpcClient, name)
	if err != nil {
		panic(err)
	}

	var target string
	if reader.URL.Path == "/" {
		target = "index.html"
	} else {
		target = reader.URL.Path[1:]
	}

	Request(writer, reader, rpcClient, addr, target)
}

func TopMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := RedirectToDefaultResourceInterceptor(
			MassaTLDInterceptor(&Interceptor{writer: w, reader: r})) //nolint:contextcheck
		if req != nil {
			handler.ServeHTTP(w, r)
		}
	})
}

type Interceptor struct {
	writer http.ResponseWriter
	reader *http.Request
}

func RedirectToDefaultResourceInterceptor(req *Interceptor) *Interceptor {
	if req == nil {
		return nil
	}

	prefixes := []string{"/browse/", "/thyra/"}

	for _, prefix := range prefixes {
		if !strings.HasPrefix(req.reader.URL.Path, prefix) {
			continue
		}

		splited := removeEmptyStrings(strings.Split(req.reader.URL.Path[len(prefix):], "/"))

		if len(splited) == 1 {
			http.Redirect(
				req.writer, req.reader,
				"http://"+req.reader.Host+prefix+splited[0]+"/"+"index.html",
				http.StatusSeeOther)

			return nil
		}
	}

	return req
}

func MassaTLDInterceptor(req *Interceptor) *Interceptor {
	if req == nil {
		return nil
	}

	massaIndex := strings.Index(req.reader.Host, ".massa")
	if massaIndex > 0 && req.reader.Host != "my.massa" {
		handleMassaDomainRequest(req.writer, req.reader, massaIndex)

		return nil
	}

	return req
}
