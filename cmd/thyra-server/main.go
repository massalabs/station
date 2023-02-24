package main

import (
	_ "embed"
	"flag"
	"log"
	"os"
	"path"

	"fyne.io/fyne/v2/app"
	"github.com/massalabs/thyra/int/api"
	"github.com/massalabs/thyra/pkg/config"
)

//go:embed version.txt
var version string

func ParseFlags() api.StartServerFlags {
	const httpPort = 80

	const httpsPort = 443

	var flags api.StartServerFlags

	configDir, _ := config.GetConfigDir()
	defaultCertFile := path.Join(configDir, "certs", "cert.pem")
	defaultCertKeyFile := path.Join(configDir, "certs", "cert-key.pem")

	flag.IntVar(&flags.Port, "http-port", httpPort, "HTTP port to listen to")
	flag.IntVar(&flags.TLSPort, "https-port", httpsPort, "HTTPS port to listen to")
	flag.StringVar(&flags.TLSCertificate, "tls-certificate", defaultCertFile, "path to certificate file")
	flag.StringVar(&flags.TLSCertificateKey, "tls-key", defaultCertKeyFile, "path to key file")
	flag.StringVar(&flags.MassaNodeServer, "node-server", "TESTNET", `Massa node that Thyra connects to. 
	Can be an IP address, a URL or one of the following values: 'TESTNET', 'LABNET', 'INNONET' or LOCALHOST`)
	flag.StringVar(&flags.DNSAddress, "dns-address", "", "Address of the DNS contract on the blockchain")
	flag.BoolVar(&flags.Version, "version", false, "Print version info and exit")

	flag.Parse()

	return flags
}

func main() {
	flags := ParseFlags()
	if flags.Version {
		log.Println("Version:", version)
		os.Exit(0)
	}

	myApp := app.New()

	go api.StartServer(&myApp, flags)

	myApp.Run()
}
