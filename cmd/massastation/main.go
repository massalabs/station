package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/massalabs/station/int/api"
	"github.com/massalabs/station/int/systray"
	"github.com/massalabs/station/int/systray/update"
	"github.com/massalabs/station/pkg/config"
)

func ParseFlags() api.StartServerFlags {
	const httpPort = 80

	const httpsPort = 443

	var flags api.StartServerFlags

	_, err := config.GetConfigDir()
	if err != nil {
		log.Fatalln("Unable to read config dir:", err,
			`MassaStation can't run without a config directory.
			Please reinstall MassaStation using the installer at https://github.com/massalabs/station and try again.`,
		)
	}

	certDir, err := config.GetCertDir()
	if err != nil {
		log.Fatalln("Unable to read cert dir:", err,
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

type logger struct {
	f   *os.File
	wrt io.Writer
}

func newLogger() *logger {
	logFilePath, err := getLogFilePath()
	if err != nil {
		log.Fatal(err)
	}

	// Create the log file if it doesn't exist
	_, err = os.Stat(logFilePath)
	if os.IsNotExist(err) {
		_, err := os.Create(logFilePath)
		if err != nil {
			log.Fatalf("error creating log file: %v", err)
		}
	}

	f, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_APPEND|os.O_TRUNC, 0o644)
	if err != nil {
		log.Fatalf("error opening log file: %v", err)
	}

	wrt := io.MultiWriter(os.Stdout, f)
	log.SetOutput(wrt)
	return &logger{f, wrt}
}

func (l *logger) close() {
	l.f.Close()
	l.wrt = nil
}

func getLogFilePath() (string, error) {
	// Get the operating system name
	osName := runtime.GOOS

	// Define the default log directory paths based on the operating system
	var logDir string
	switch osName {
	case "linux":
		logDir = "/usr/local/share/massastation"
	case "darwin":
		logDir = "/Library/Logs/"
	case "windows":
		logDir = os.Getenv("USERPROFILE")
	default:
		return "", fmt.Errorf("unsupported operating system: %s", osName)
	}

	// Create the logs directory if it doesn't exist
	err := os.MkdirAll(logDir, os.ModePerm)
	if err != nil {
		return "", err
	}

	// Return the log file path
	return filepath.Join(logDir, "massastation.log"), nil
}

func main() {
	logger := newLogger()
	defer logger.close()
	flags := ParseFlags()
	if flags.Version {
		log.Println("Version:", config.Version)
		os.Exit(0)
	}

	networkManager, err := config.NewNetworkManager()
	if err != nil {
		log.Fatal("Failed to create NetworkManager:", err)
	}

	stationGUI, systrayMenu := systray.MakeGUI()
	server := api.NewServer(flags)

	update.StartUpdateCheck(&stationGUI, systrayMenu)

	stationGUI.Lifecycle().SetOnStopped(server.Stop)
	stationGUI.Lifecycle().SetOnStarted(func() {
		server.Start(networkManager)
	})

	stationGUI.Run()
}
