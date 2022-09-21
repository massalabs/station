package resource

// Front HTML

//go:generate textFileToGoConst -in html/front/wallet.css -o ../pkg/front/wallet/css.go -p wallet -c CSS
//go:generate textFileToGoConst -in html/front/wallet.html -o ../pkg/front/wallet/html.go -p wallet -c HTML
//go:generate textFileToGoConst -in html/front/wallet.js -o ../pkg/front/wallet/js.go -p wallet -c JS
//go:generate textFileToGoConst -in html/front/website.css -o ../pkg/front/website/css.go -p website -c CSS
//go:generate textFileToGoConst -in html/front/website.html -o ../pkg/front/website/html.go -p website -c HTML
//go:generate textFileToGoConst -in html/front/website.js -o ../pkg/front/website/js.go -p website -c JS
//go:generate textFileToGoConst -in html/front/errors.js -o ../pkg/front/errors.go -p front -c Errors
//go:generate textFileToGoConst -in html/front/common.js -o ../pkg/front/common.go -p front -c Common
//go:generate textFileToGoConst -in html/front/logo_banner.webp -o ../pkg/front/logoBanner.go -p front -c LogoBanner
//go:generate textFileToGoConst -in html/front/logo.png -o ../pkg/front/logo.go -p front -c Logo

// API server certificate

//go:generate textFileToGoConst -in api/certificate/unsecure.crt -o ../api/swagger/server/certificate.go -p server -c UnsecureCertificate
//go:generate textFileToGoConst -in api/certificate/unsecure.key -o ../api/swagger/server/key.go -p server -c UnsecureKey

// Contracts

//go:generate textFileToGoConst -in sc/websiteStorer.wasm -o ../pkg/sc/websiteStorer.go -p sc -c WebsiteStorer

// API swagger

//go:generate swagger generate server --quiet --target ../api/swagger/server --name thyra-server --spec api/swagger.yml --exclude-main
