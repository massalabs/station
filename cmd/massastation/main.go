package main

import (
	_ "embed"
	"flag"
	"log"
	"net/url"
	"os"
	"path"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	//nolint:typecheck
	"fyne.io/fyne/v2/driver/desktop"

	"github.com/massalabs/thyra/int/api"
	"github.com/massalabs/thyra/pkg/config"
)

//go:embed logo.png
var logo []byte

func openURL(app *fyne.App, urlToOpen string) {
	u, err := url.Parse(urlToOpen)
	if err != nil {
		log.Fatal(err)
	}

	err = (*app).OpenURL(u)
	if err != nil {
		log.Fatal(err)
	}
}

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
	Can be an IP address, a URL or one of the following values: 'TESTNET', 'LABNET', 'INNONET' or LOCALHOST`)
	flag.StringVar(&flags.DNSAddress, "dns-address", "", "Address of the DNS contract on the blockchain")
	flag.BoolVar(&flags.Version, "version", false, "Print version info and exit")

	flag.Parse()

	return flags
}

func makeGUI() (fyne.App, *fyne.Menu) {
	stationGUI := app.New()
	menu := fyne.NewMenu("Thyra Desktop")
	menu.Refresh()

	if desk, ok := stationGUI.(desktop.App); ok {
		icon := fyne.NewStaticResource("logo", logo)
		titleMenu := fyne.NewMenuItem("MassaStation", nil)
		homeShortCutMenu := fyne.NewMenuItem("Open MassaStation", nil)
		testMenu := fyne.NewMenuItem("Test", nil)

		titleMenu.Disabled = true
		titleMenu.Icon = icon

		testMenu.Action = func() {
			notification := fyne.NewNotification("MassaStation", "MassaStation is running in the background")
			stationGUI.SendNotification(notification)
		}

		homeShortCutMenu.Action = func() {
			openURL(&stationGUI, "http://my.massa/thyra/home")
		}

		menu.Items = append(menu.Items,
			titleMenu,
			fyne.NewMenuItemSeparator(),
			homeShortCutMenu,
			// testMenu,
			fyne.NewMenuItemSeparator(),
		)

		desk.SetSystemTrayIcon(icon)
		desk.SetSystemTrayMenu(menu)
	}

	return stationGUI, menu
}

func main() {
	flags := ParseFlags()
	if flags.Version {
		log.Println("Version:", config.Version)
		os.Exit(0)
	}

	stationGUI, _ := makeGUI()
	server := api.NewServer(flags)

	stationGUI.Lifecycle().SetOnStopped(server.Stop)
	stationGUI.Lifecycle().SetOnStarted(server.Start)

	stationGUI.Run()
}
