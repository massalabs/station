import ctypes
import logging
import os
import platform
import zipfile

from installer import Installer

class WindowsInstaller(Installer):
    ACRYLIC_DNS_PROXY_URL = "https://sourceforge.net/projects/acrylic/files/Acrylic/2.1.1/Acrylic-Portable.zip/download"
    ACRYLIC_DNS_PROXY_FILENAME = "Acrylic-Portable.zip"
    ACRYLIC_HOST_FILE = "AcrylicHosts.txt"
    DEFAULT_ACRYLIC_PATH = "C:\Program Files (x86)\Acrylic DNS Proxy"

    def __init__(self):
        super().__init__()
        self.THYRA_SERVER_FILENAME = "thyra-server.exe"
        self.THYRA_APP_FILENAME = "thyra-app.exe"
        self.MKCERT_FILENAME = "mkcert.exe"

        if platform.machine() == "AMD64":
            self.THYRA_SERVER_URL = "https://github.com/massalabs/thyra/releases/latest/download/thyra-server_windows_amd64"
            self.THYRA_APP_URL = "https://github.com/massalabs/Thyra-Menu-Bar-App/releases/latest/download/ThyraApp_windows-amd64.exe"
            self.MKCERT_URL = "https://dl.filippo.io/mkcert/latest?for=windows/amd64"
            self.THYRA_WALLET_PLUGIN_URL = "https://github.com/massalabs/thyra-plugin-wallet/releases/download/v0.0.10/wallet-plugin_windows-amd64.zip"
            self.THYRA_WALLET_BINARY_FILENAME = "wallet-plugin_windows-amd64.exe"
            self.THYRA_WALLET_ZIP_FILENAME = "wallet-plugin_windows-amd64.zip"
        else:
            self.printErrorAndExit(f"Unsupported architecture {platform.machine()}")

    @staticmethod
    def printErrorAndExit(error):
        logging.error(error)
        os.system("pause")
        os._exit(-1)

    def unzipAcrylic(self):
        logging.info("Unzipping Acrylic...")
        try:
            with zipfile.ZipFile(self.ACRYLIC_DNS_PROXY_FILENAME, 'r') as zip_ref:
                zip_ref.extractall(self.DEFAULT_ACRYLIC_PATH)
        except:
            self.printErrorAndExit("Could not unzip Acrylic")

    def configureAcrylic(self):
        logging.info("Configuring Acrylic...")
        f = open(os.path.join(self.DEFAULT_ACRYLIC_PATH, self.ACRYLIC_HOST_FILE), "r+")
        if f.read().find("127.0.0.1 *.massa") != -1:
            f.close()
            return
        f.write("\n127.0.0.1 *.massa")
        f.close()

        logging.info("Restarting Acrylic service...")
        self.executeCommand("NET STOP AcrylicDNSProxySvc & NET START AcrylicDNSProxySvc", True)

    def configureNetworkInterface(self):
        logging.info("Configuring network interface...")
        stdout, _stderr = self.executeCommand("wmic nic where \"netenabled=true\" get netconnectionID", True)

        networkAdaptersNames = list(filter(None, stdout.split('\n')))
        networkAdaptersNames.pop(0)
        networkAdaptersNames = [name.strip() for name in networkAdaptersNames]

        for name in networkAdaptersNames:
            self.executeCommand(f'NETSH interface ipv4 set dnsservers "{name}" static 127.0.0.1 primary')
            self.executeCommand(f'NETSH interface ipv6 set dnsservers "{name}" static ::1 primary')

        logging.info("Network interface configured")

    def setupDNS(self):
        logging.info("Setting up DNS...")

        if os.path.exists(self.DEFAULT_ACRYLIC_PATH):
            logging.info("Acrylic DNS Proxy is already installed")
        else:
            cwd = os.getcwd()

            logging.info("Installing Acrylic DNS Proxy...")
            self.downloadFile(self.ACRYLIC_DNS_PROXY_URL, self.ACRYLIC_DNS_PROXY_FILENAME)
            self.unzipAcrylic()

            os.chdir(self.DEFAULT_ACRYLIC_PATH)
            self.executeCommand(os.path.join(self.DEFAULT_ACRYLIC_PATH, "InstallAcrylicService.bat"))
            self.configureNetworkInterface()
            os.chdir(cwd)

            try:
                os.remove(self.ACRYLIC_DNS_PROXY_FILENAME)
            except:
                logging.error("Could not remove Acrylic DNS Proxy installer")

        self.configureAcrylic()


if __name__ == "__main__":
    if platform.system() != "Windows":
        WindowsInstaller.printErrorAndExit("This script is only compatible with Windows")

    if ctypes.windll.shell32.IsUserAnAdmin() == 0:
        WindowsInstaller.printErrorAndExit("This script must be run as administrator")

    WindowsInstaller().startInstall()
    os.system("pause")
