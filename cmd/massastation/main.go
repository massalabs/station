package main

import (
	"flag"
	"os"
	"path"

	"github.com/massalabs/station/int/api"
	"github.com/massalabs/station/int/systray"
	"github.com/massalabs/station/int/systray/update"
	"github.com/massalabs/station/pkg/config"
	"github.com/massalabs/station/pkg/dirs"
	"github.com/massalabs/station/pkg/logger"
	"github.com/massalabs/station/pkg/plugin"
)

func ParseFlags() api.StartServerFlags {
	const httpPort = 80

	const httpsPort = 443

	var flags api.StartServerFlags

	_, err := dirs.GetConfigDir()
	if err != nil {
		logger.Logger.Error(
			"Unable to read config dir: %s\n%s",
			err,
			`MassaStation can't run without a config directory.\n
			Please reinstall MassaStation using the installer at https://github.com/massalabs/station and try again.`)
	}

	certDir, err := dirs.GetCertDir()
	if err != nil {
		logger.Logger.Fatal("Unable to read cert dir:%s\n%s", err,
			`MassaStation can't run without a certificate directory.
			Please reinstall MassaStation using the installer at https://github.com/massalabs/station and try again.`,
		)
	}

	defaultCertFile := path.Join(certDir, "cert.pem")
	defaultCertKeyFile := path.Join(certDir, "cert-key.pem")

	flag.IntVar(&flags.Port, "http-port", httpPort, "HTTP port to listen to")
	flag.IntVar(&flags.TLSPort, "https-port", httpsPort, "HTTPS port to listen to")
	flag.StringVar(&flags.TLSCertificate, "tls-certificate", defaultCertFile, "path to certificate file")
	flag.StringVar(&flags.TLSCertificateKey, "tls-key", defaultCertKeyFile, "path to key file")
	flag.StringVar(&flags.MassaNodeServer, "node-server", "TESTNET", `Massa node that MassaStation connects to. 
	Can be an IP address, a URL or one of the following values: 'TESTNET', 'LABNET', 'BUILDNET' or LOCALHOST`)
	flag.StringVar(&flags.DNSAddress, "dns-address", "", "Address of the DNS contract on the blockchain")
	flag.BoolVar(&flags.Version, "version", false, "Print version info and exit")

	flag.Parse()

	return flags
}

func main() {
	logger.Logger = logger.NewLogger()
	defer logger.Logger.Sync()

	flags := ParseFlags()
	if flags.Version {
		logger.Logger.Infof("Version:%s", config.Version)
		os.Exit(0)
	}

	err := config.Check()
	if err != nil {
		logger.Logger.Fatalf("Error with you current system configuration: %s", err.Error())
	}

	networkManager, err := config.NewNetworkManager()
	if err != nil {
		logger.Logger.Fatalf("Failed to create NetworkManager:%s", err.Error())
	}

	pluginManager := plugin.NewManager()

	stationGUI, systrayMenu := systray.MakeGUI()
	server := api.NewServer(flags)

	update.StartUpdateCheck(&stationGUI, systrayMenu)

	stationGUI.Lifecycle().SetOnStopped(func() {
		pluginManager.Stop()
		server.Stop()
	})
	stationGUI.Lifecycle().SetOnStarted(func() {
		server.Start(networkManager, pluginManager)
		err := pluginManager.RunAll()
		if err != nil {
			logger.Logger.Fatalf("while running all plugins: %w", err.Error())
		}
	})

	stationGUI.Run()
}
