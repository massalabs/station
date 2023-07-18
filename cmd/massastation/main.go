package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/massalabs/station/int/api"
	"github.com/massalabs/station/int/configuration"
	"github.com/massalabs/station/int/initialize"
	"github.com/massalabs/station/int/systray"
	"github.com/massalabs/station/int/systray/update"
	"github.com/massalabs/station/pkg/config"
	"github.com/massalabs/station/pkg/logger"
	"github.com/massalabs/station/pkg/plugin"
)

func ParseFlags() api.StartServerFlags {
	const httpPort = 80

	const httpsPort = 443

	var flags api.StartServerFlags

	_, err := configuration.Path()
	if err != nil {
		logger.Error(
			"Unable to read config dir: %s\n%s",
			err,
			`MassaStation can't run without a config directory.\n
			Please reinstall MassaStation using the installer at https://github.com/massalabs/station and try again.`)
	}

	flag.IntVar(&flags.Port, "http-port", httpPort, "HTTP port to listen to")
	flag.IntVar(&flags.TLSPort, "https-port", httpsPort, "HTTPS port to listen to")
	flag.StringVar(&flags.MassaNodeServer, "node-server", "TESTNET", `Massa node that MassaStation connects to. 
	Can be an IP address, a URL or one of the following values: 'TESTNET', 'LABNET', 'BUILDNET' or LOCALHOST`)
	flag.StringVar(&flags.DNSAddress, "dns-address", "", "Address of the DNS contract on the blockchain")
	flag.BoolVar(&flags.Version, "version", false, "Print version info and exit")

	flag.Parse()

	return flags
}

func main() {
	flags := ParseFlags()
	if flags.Version {
		//nolint:forbidigo
		fmt.Printf("Version:%s\n", config.Version)
		os.Exit(0)
	}

	err := initialize.Logger()
	if err != nil {
		log.Fatalf("while initializing logger: %s", err.Error())
	}

	defer logger.Close()

	err = config.Check()
	if err != nil {
		logger.Fatalf("Error with you current system configuration: %s", err.Error())
	}

	networkManager, err := config.NewNetworkManager()
	if err != nil {
		logger.Fatalf("Failed to create NetworkManager:%s", err.Error())
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
			logger.Fatalf("while running all plugins: %w", err.Error())
		}
	})

	stationGUI.Run()
}
