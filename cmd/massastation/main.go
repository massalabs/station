package main

import (
	"flag"
	"log"
	"os"
	"path"

	"github.com/massalabs/thyra/int/api"
	"github.com/massalabs/thyra/int/systray"
	"github.com/massalabs/thyra/pkg/config"
)

func ParseFlags() api.StartServerFlags {
	const httpPort = 80

	const httpsPort = 443

	var flags api.StartServerFlags

	_, err := config.GetConfigDir()
	if err != nil {
		log.Fatalln("Unable to read config dir:", err,
			`MassaStation can't run without a config directory.
			Please reinstall MassaStation using the installer at https://github.com/massalabs/thyra and try again.`,
		)
	}

	certDir, err := config.GetCertDir()
	if err != nil {
		log.Fatalln("Unable to read cert dir:", err,
			`MassaStation can't run without a certificate directory.
			Please reinstall MassaStation using the installer at https://github.com/massalabs/thyra and try again.`,
		)
	}

	defaultCertFile := path.Join(certDir, "cert.pem")
	defaultCertKeyFile := path.Join(certDir, "cert-key.pem")

	flag.IntVar(&flags.Port, "http-port", httpPort, "HTTP port to listen to")
	flag.IntVar(&flags.TLSPort, "https-port", httpsPort, "HTTPS port to listen to")
	flag.StringVar(&flags.TLSCertificate, "tls-certificate", defaultCertFile, "path to certificate file")
	flag.StringVar(&flags.TLSCertificateKey, "tls-key", defaultCertKeyFile, "path to key file")
	flag.StringVar(&flags.MassaNodeServer, "node-server", "TESTNET", `Massa node that Thyra connects to. 
	Can be an IP address, a URL or one of the following values: 'TESTNET', 'LABNET', 'BUILDNET' or LOCALHOST`)
	flag.StringVar(&flags.DNSAddress, "dns-address", "", "Address of the DNS contract on the blockchain")
	flag.BoolVar(&flags.Version, "version", false, "Print version info and exit")

	flag.Parse()

	return flags
}

func main() {
	flags := ParseFlags()
	if flags.Version {
		log.Println("Version:", config.Version)
		os.Exit(0)
	}

	stationGUI, systrayMenu := systray.MakeGUI()
	server := api.NewServer(flags)

	stationGUI.Lifecycle().SetOnStopped(server.Stop)
	stationGUI.Lifecycle().SetOnStarted(server.Start)

	stationGUI.Run()
}
