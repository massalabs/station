import ctypes
import logging
import os
import platform
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
        else:
            self.printErrorAndExit(f"Unsupported architecture {platform.machine()}")

    def printErrorAndExit(self, error):
        logging.error(error)
        os.system("pause")
        os._exit(-1)

    def configureAchylic(self):
        logging.info("Configuring Acrylic...")
        f = open(self.DEFAULT_ACRYLIC_PATH + "\\" + self.ACRYLIC_HOST_FILE, "r+")
        if f.read().find("127.0.0.1 *.massa") != -1:
            f.close()
            return
        f.write("\n127.0.0.1 *.massa")
        f.close()

        self.executeCommand("NET STOP AcrylicDNSProxySvc", True)
        self.executeCommand("NET START AcrylicDNSProxySvc", True)

    def configureNetworkInterface(self):
        commandOutput = self.executeCommand("wmic nic where \"netenabled=true\" get netconnectionID", True)

        networkAdpatersNames = list(filter(None, commandOutput.split('\n')))
        networkAdpatersNames.pop(0)
        networkAdpatersNames = [name.strip() for name in networkAdpatersNames]

        for name in networkAdpatersNames:
            self.executeCommand("NETSH interface ipv4 set dnsservers " + name + " static 127.0.0.1 primary", True)
            self.executeCommand("NETSH interface ipv6 set dnsservers " + name + " static ::1 primary", True)

    def setupDNS(self):
        logging.info("Setting up DNS...")

        if os.path.exists(self.DEFAULT_ACRYLIC_PATH):
            logging.info("Acrylic DNS Proxy is already installed")
        else:
            logging.info("Installing Acrylic DNS Proxy...")
            self.downloadFile(self.ACRYLIC_DNS_PROXY_URL, self.ACRYLIC_DNS_PROXY_FILENAME)
            self.executeCommand(f"Expand-Archive {self.ACRYLIC_DNS_PROXY_FILENAME} -DestinationPath {self.DEFAULT_ACRYLIC_PATH}", True)
            self.executeCommand(f"{self.DEFAULT_ACRYLIC_PATH}\InstallAcrylicService.bat", True)
            self.configureNetworkInterface()

            try:
                os.remove(self.ACRYLIC_DNS_PROXY_FILENAME)
            except:
                logging.error("Could not remove Acrylic DNS Proxy installer")

        self.configureAchylic()


if __name__ == "__main__":
    if platform.system() != "Windows":
        logging.error("This script is only compatible with Windows")
        os._exit(-1)

    if ctypes.windll.shell32.IsUserAnAdmin() == 0:
        logging.error("This script must be run as administrator")
        os._exit(-1)

    WindowsInstaller().startInstall()
    os.system("pause")
