package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/massalabs/station/int/api"
	"github.com/massalabs/station/int/config"
	"github.com/massalabs/station/int/initialize"
	"github.com/massalabs/station/int/systray"
	"github.com/massalabs/station/int/systray/update"
	"github.com/massalabs/station/pkg/dirs"
	"github.com/massalabs/station/pkg/logger"
	"github.com/massalabs/station/pkg/plugin"
)

type StartFlags struct {
	Version bool
	Repair  bool
}

func ParseFlags() (api.StartServerFlags, StartFlags) {
	const httpPort = 80

	const httpsPort = 443

	var serverFlags api.StartServerFlags

	var startFlags StartFlags

	_, err := dirs.GetConfigDir()
	if err != nil {
		logger.Errorf(
			"Unable to read config dir: %s\n%s",
			err,
			`MassaStation can't run without a config directory.\n
			Please reinstall MassaStation using the installer at https://github.com/massalabs/station and try again.`)
	}

	flag.IntVar(&serverFlags.Port, "http-port", httpPort, "HTTP port to listen to")
	flag.IntVar(&serverFlags.TLSPort, "https-port", httpsPort, "HTTPS port to listen to")
	flag.BoolVar(&startFlags.Version, "version", false, "Print version info and exit")
	flag.BoolVar(&startFlags.Repair, "repair", false, "Repair MassaStation")

	flag.Parse()

	return serverFlags, startFlags
}

func main() {
	err := initialize.Logger()
	if err != nil {
		log.Fatalf("while initializing logger: %s", err.Error())
	}

	defer logger.Close()

	serverFlags, startFlags := ParseFlags()
	if startFlags.Version {
		//nolint:forbidigo
		fmt.Printf("Version:%s\n", config.Version)
		logger.Close()
		//nolint:gocritic
		os.Exit(0)
	}

	err = config.Check()
	if err != nil {
		logger.Fatalf("Error with you current system configuration: %s", err.Error())
	}

	if startFlags.Repair {
		os.Exit(0)
	}

	networkManager, err := config.NewNetworkManager()
	if err != nil {
		logger.Fatalf("Failed to create NetworkManager:%s", err.Error())
	}

	pluginManager := plugin.NewManager()

	stationGUI, systrayMenu := systray.MakeGUI()
	server := api.NewServer(serverFlags)

	update.StartUpdateCheck(&stationGUI, systrayMenu)

	stationGUI.Lifecycle().SetOnStopped(func() {
		pluginManager.StopAll()
		server.Stop()
	})
	stationGUI.Lifecycle().SetOnStarted(func() {
		server.Start(networkManager, pluginManager)
		err := pluginManager.RunAll()
		if err != nil {
			logger.Fatalf("while running all plugins: %w", err)
		}
	})

	stationGUI.Run()
}
