import logging
import platform

from installer import Installer

class MacOSInstaller(Installer):
    def __init__(self):
        super().__init__()
        self.THYRA_SERVER_FILENAME = "thyra_server"
        self.THYRA_APP_FILENAME = "thyra-app"
        self.MKCERT_FILENAME = "mkcert"

        if platform.machine() == "arm64":
            self.THYRA_SERVER_URL = "https://github.com/massalabs/thyra/releases/latest/download/thyra-server_darwin_arm64"
            self.THYRA_APP_URL = "https://github.com/massalabs/Thyra-Menu-Bar-App/releases/latest/download/ThyraApp_darwin-arm64"
            self.MKCERT_URL = "https://dl.filippo.io/mkcert/latest?for=darwin/arm64"
        elif platform.machine() == "x86_64":
            self.THYRA_SERVER_URL = "https://github.com/massalabs/thyra/releases/latest/download/thyra-server_darwin_amd64"
            self.THYRA_APP_URL = "https://github.com/massalabs/Thyra-Menu-Bar-App/releases/latest/download/ThyraApp_darwin-amd64"
            self.MKCERT_URL = "https://dl.filippo.io/mkcert/latest?for=darwin/amd64"
        else:
            self.printErrorAndExit(f"Unsupported architecture {platform.machine()}")

    def configureDNSMasq(self):
        logging.info("Configuring DNSMasq...")

        self.executeCommand("sudo bash -c 'echo ""address=/.massa/127.0.0.1"" > $(brew --prefix)/etc/dnsmasq.d/massa.conf'", True)
        self.executeCommand("sudo mkdir -p /etc/resolver", True)
        self.executeCommand("sudo bash -c 'echo ""nameserver 127.0.0.1"" > /etc/resolver/massa'", True)

        logging.info("Restarting DNSMasq...")
        self.executeCommand("sudo brew services restart dnsmasq", True)

    def setupDNS(self):
        stdout, _stderr = self.executeCommand("sudo lsof -i :53", True, allow_failure=True)        
        runningDNS = ""
        if stdout:
            runningDNS = stdout.splitlines()[1].split()
            runningDNS = runningDNS[:-1]

        if runningDNS == "dnsmasq":
            logging.info("dnsmasq is already installed")
            self.configureDNSMasq()
        elif runningDNS == "":
            logging.info("Installing dnsmasq...")
            self.executeCommand("brew install dnsmasq", True)
            self.configureDNSMasq()
        else:
            logging.info(runningDNS)
            self.printErrorAndExit("Unsupported DNS server")


if __name__ == "__main__":
    if platform.system() != "Darwin":
        MacOSInstaller.printErrorAndExit("This script is only for MacOS")

    MacOSInstaller().startInstall()
