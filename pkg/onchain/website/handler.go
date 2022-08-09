package website

import (
	"net/http"
	"strings"

	"github.com/massalabs/thyra/pkg/front"
	fwallet "github.com/massalabs/thyra/pkg/front/wallet"
	"github.com/massalabs/thyra/pkg/front/website"
	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/onchain/dns"
)

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
		c := node.NewDefaultClient()
		Request(w, r, c, splited[0], splited[1])
	default:
		pathNotFound(w, r)
	}
}

func handleMassaDomainRequest(w http.ResponseWriter, r *http.Request, index int) {
	name := r.Host[:index]

	rpcClient := node.NewDefaultClient()

	addr, err := dns.Resolve(rpcClient, name)
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

//TO REWORK
func HandlerFunc(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		massaTLD := strings.Index(r.Host, ".massa")
		if strings.HasPrefix(r.Host, "webuploader.mythyra.massa") && strings.Index(r.URL.Path, ".") != -1 {
			HandleWebsiteUploaderManagementRequest(w, r)
		} else if strings.HasPrefix(r.Host, "wallet.mythyra.massa") && strings.Index(r.URL.Path, ".") != -1 {
			HandleWalletManagementRequest(w, r)
		} else if massaTLD > 0 && strings.Index(r.Host, "mythyra") == -1 {
			handleMassaDomainRequest(w, r, massaTLD)
		} else if strings.HasPrefix(r.URL.Path, "/website") {
			handleAPIRequest(w, r)
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
	case "logo_banner.webp":
		fileText = front.LogoBanner
	case "logo.png":
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
	case "wallet.css":
		fileText = fwallet.CSS
	case "wallet.js":
		fileText = fwallet.JS
	case "wallet.html":
		fileText = fwallet.HTML
	case "logo_banner.webp":
		fileText = front.LogoBanner
	case "logo.png":
		fileText = front.Logo
	}
	setContentType(resource, w)

	_, err := w.Write([]byte(fileText))
	if err != nil {
		panic(err)
	}
}
