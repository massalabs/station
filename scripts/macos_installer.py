import logging
import os
import platform

from installer import Installer

class MacOSInstaller(Installer):
    def __init__(self):
        super().__init__()
        self.MASSASTATION_SERVER_FILENAME = "thyra-server"
        self.MASSASTATION_APP_FILENAME = "thyra-app"

        self.MASSASTATION_CONFIG_FOLDER_PATH = "/usr/local/share/massastation"
        self.MASSASTATION_PLUGINS_PATH = "/usr/local/share/massastation/plugins"
        self.CERTIFICATIONS_FOLDER_PATH = "/etc/massastation/certs"
        self.MKCERT_FILENAME = "mkcert"

        self.SUDO_INSTALLATION = True

        if platform.machine() == "arm64":
            self.MASSASTATION_SERVER_URL = "https://github.com/massalabs/thyra/releases/latest/download/thyra-server_darwin_arm64"
            self.MASSASTATION_APP_URL = "https://github.com/massalabs/Thyra-Menu-Bar-App/releases/latest/download/ThyraApp_darwin-arm64"
            self.MKCERT_URL = "https://dl.filippo.io/mkcert/latest?for=darwin/arm64"
            self.MASSASTATION_WALLET_PLUGIN_URL = "https://github.com/massalabs/thyra-plugin-wallet/releases/latest/download/wallet-plugin_darwin-arm64.zip"
            self.MASSASTATION_WALLET_BINARY_FILENAME = "wallet-plugin_darwin-arm64"
            self.MASSASTATION_WALLET_ZIP_FILENAME = "wallet-plugin_darwin-arm64.zip"
        elif platform.machine() == "x86_64":
            self.MASSASTATION_SERVER_URL = "https://github.com/massalabs/thyra/releases/latest/download/thyra-server_darwin_amd64"
            self.MASSASTATION_APP_URL = "https://github.com/massalabs/Thyra-Menu-Bar-App/releases/latest/download/ThyraApp_darwin-amd64"
            self.MKCERT_URL = "https://dl.filippo.io/mkcert/latest?for=darwin/amd64"
            self.MASSASTATION_WALLET_PLUGIN_URL = "https://github.com/massalabs/thyra-plugin-wallet/releases/latest/download/wallet-plugin_darwin-amd64.zip"
            self.MASSASTATION_WALLET_BINARY_FILENAME = "wallet-plugin_darwin-amd64"
            self.MASSASTATION_WALLET_ZIP_FILENAME = "wallet-plugin_darwin-amd64.zip"
        else:
            self.printErrorAndExit(f"Unsupported architecture {platform.machine()}")

    def configureNetworkInterface(self):
        logging.info("Configuring network interface...")
        stdout, _stderr = self.executeCommand("networksetup -listallnetworkservices", True)

        networkAdaptersNames = list(filter(None, stdout.split('\n')))
        networkAdaptersNames.pop(0)
        networkAdaptersNames = [adapter.strip() for adapter in networkAdaptersNames]

        for adapter in networkAdaptersNames:
            self.executeCommand(f'networksetup -setdnsservers "{adapter}" 127.0.0.1', True)

        logging.info("Network interface configured")

    def configureDNSMasq(self):
        logging.info("Configuring DNSMasq...")

        servers = [
            "8.8.8.8",
            "8.8.4.4",
        ]

        self.executeCommand("sudo bash -c 'echo -e ""address=/.massa/127.0.0.1"" > $(brew --prefix)/etc/dnsmasq.d/massa.conf'", True)
        self.executeCommand("sudo bash -c 'echo -e ""no-resolv"" >> $(brew --prefix)/etc/dnsmasq.d/massa.conf'", True)
        for server in servers:
            test = f'echo -e ""server={server}"" >> $(brew --prefix)/etc/dnsmasq.d/massa.conf'
            self.executeCommand(f"sudo bash -c '{test}'", True)

        self.executeCommand("sudo mkdir -p /etc/resolver", True)
        self.executeCommand("sudo bash -c 'echo ""nameserver 127.0.0.1"" > /etc/resolver/massa'", True)

        logging.info("Restarting DNSMasq...")
        self.executeCommand("sudo brew services restart dnsmasq", True)

    def setupDNS(self):
        stdout, _ = self.executeCommand("sudo lsof -i :53", True, allow_failure=True)        
        runningDNS = ""
        if stdout:
            runningDNS = stdout.splitlines()[1].split()
            runningDNS = runningDNS[0]

        if runningDNS == "dnsmasq":
            logging.info("dnsmasq is already installed")
        elif runningDNS == "":
            logging.info("Installing dnsmasq...")
            self.executeCommand("brew install dnsmasq", True)
        else:
            logging.info(runningDNS)
            self.printErrorAndExit(f"Unsupported DNS server {runningDNS}")
        self.configureDNSMasq()
        self.configureNetworkInterface()

    def generateCACertificate(self):
        stdout, _ = self.executeCommand("find /Applications/ -type d -iname '*Firefox*.app'", True, allow_failure=True)
        if stdout and "Firefox" in stdout:
            self.executeCommand("brew install nss", True)

        if not os.path.exists(self.CERTIFICATIONS_FOLDER_PATH):
            self.executeCommand(f"sudo mkdir -p {self.CERTIFICATIONS_FOLDER_PATH}", True)

        super().generateCACertificate()

if __name__ == "__main__":
    if platform.system() != "Darwin":
        MacOSInstaller.printErrorAndExit("This script is only for MacOS")

    MacOSInstaller().startInstall()
