package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/massalabs/station/int/api"
	"github.com/massalabs/station/int/config"
	"github.com/massalabs/station/int/configuration"
	"github.com/massalabs/station/int/initialize"
	"github.com/massalabs/station/int/systray"
	"github.com/massalabs/station/int/systray/update"
	"github.com/massalabs/station/int/systray/utils"
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

	flag.IntVar(&serverFlags.Port, "http-port", httpPort, "HTTP port to listen to")
	flag.IntVar(&serverFlags.TLSPort, "https-port", httpsPort, "HTTPS port to listen to")
	flag.BoolVar(&startFlags.Version, "version", false, "Print version info and exit")
	flag.BoolVar(&startFlags.Repair, "repair", false, "Repair MassaStation")

	flag.Parse()

	return serverFlags, startFlags
}

//nolint:funlen
func main() {
	serverFlags, stationStartFlags := ParseFlags()

	configDir, err := configuration.Path()
	if err != nil {
		log.Fatalf(
			"Unable to read config dir: %s\n%s",
			err,
			`MassaStation can't run without a config directory.\n
			Please reinstall MassaStation using the installer at https://github.com/massalabs/station and try again.`)
	}

	err = initialize.Logger(stationStartFlags.Repair, configDir)
	if err != nil {
		log.Fatalf("while initializing logger: %s", err.Error())
	}

	defer logger.Close()

	if stationStartFlags.Version {
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

	if stationStartFlags.Repair {
		logger.Infof("Repair process completed. Please check the logs for any potential errors.")
		os.Exit(0)
	}

	networkManager, err := config.NewNetworkManager()
	if err != nil {
		logger.Fatalf("Failed to create NetworkManager: %s", err.Error())
	}

	pluginManager := plugin.NewManager(configDir)

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
		utils.OpenURL(&stationGUI, fmt.Sprintf("https://%s", config.MassaStationURL))
	})

	stationGUI.Run()
}
