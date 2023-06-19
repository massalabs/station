import logging
import platform

from installer import Installer

class LinuxInstaller(Installer):
    def __init__(self):
        super().__init__()

        self.MASSASTATION_INSTALL_FOLDER_PATH = "/usr/local/bin"
        self.MASSASTATION_SERVER_FILENAME = "thyra-server"
        self.MASSASTATION_APP_FILENAME = "thyra-app"

        self.MASSASTATION_CONFIG_FOLDER_PATH = "/usr/local/share/massastation"
        self.MASSASTATION_PLUGINS_PATH = "/usr/local/share/massastation/plugins"
        self.CERTIFICATIONS_FOLDER_PATH = "/etc/massastation/certs"
        self.MKCERT_FILENAME = "mkcert"

        self.MASSASTATION_WALLET_PLUGIN_URL = "https://github.com/massalabs/thyra-plugin-wallet/releases/latest/download/wallet-plugin_linux-amd64.zip"
        self.MASSASTATION_WALLET_BINARY_FILENAME = "wallet-plugin_linux-amd64"
        self.MASSASTATION_WALLET_ZIP_FILENAME = "wallet-plugin_linux-amd64.zip"

        self.MASSASTATION_SERVER_URL = "https://github.com/massalabs/station/releases/latest/download/thyra-server_linux_amd64"
        self.MASSASTATION_APP_URL = "https://github.com/massalabs/Thyra-Menu-Bar-App/releases/latest/download/ThyraApp_linux-amd64"
        self.MKCERT_URL = "https://dl.filippo.io/mkcert/latest?for=linux/amd64"

        self.SUDO_INSTALLATION = True

    def configureNetworkManager(self):
        logging.info("Configuring NetworkManager...")
        self.executeCommand("sudo cp /etc/NetworkManager/NetworkManager.conf /etc/NetworkManager/NetworkManager.conf.backup_thyra_install", True)
        stdout, _stderr = self.executeCommand("grep '^dns=' /etc/NetworkManager/NetworkManager.conf", True, allow_failure=True)
        dns = stdout.removeprefix("dns=").strip()
        if dns == "dnsmasq":
            logging.info("NetworkManager is already configured")
        elif dns == "":
            # add dns=dnsmasq after [main] line
            self.executeCommand("sudo sed -i 's/^\[main\]$/\[main\]\\ndns=dnsmasq/' /etc/NetworkManager/NetworkManager.conf", True)
        else:
            self.executeCommand("sudo sed -i 's/^dns=.*/dns=dnsmasq/g' /etc/NetworkManager/NetworkManager.conf", True)

    def configureDNSMasq(self):
        logging.info("Configuring DNSMasq...")

        self.executeCommand("sudo mkdir -p /etc/NetworkManager/dnsmasq.d", True)
        self.executeCommand("sudo bash -c 'echo 'address=/.massa/127.0.0.1' > /etc/NetworkManager/dnsmasq.d/massa.conf'", True)
        self.executeCommand("sudo mv /etc/resolv.conf /etc/resolv.conf.backup_thyra_install", True)
        self.executeCommand("sudo ln -s /run/NetworkManager/resolv.conf /etc/resolv.conf", True)
        self.executeCommand("sudo systemctl restart NetworkManager", True)

    def setupDNS(self):
        stdout, _stderr = self.executeCommand("sudo lsof -i :53", True, allow_failure=True)
        runningDNS = ""
        if stdout:
            runningDNS = stdout.splitlines()[1].split()
            runningDNS = runningDNS[0]

        if runningDNS == "":
            logging.info("Installing dnsmasq...")
            self.executeCommand("sudo apt-get install -y dnsmasq", True)
            self.configureNetworkManager()
            self.configureDNSMasq()
        elif runningDNS == "dnsmasq":
            logging.info("dnsmasq is already installed")
            self.configureDNSMasq()
        elif runningDNS == "systemd-r":
            logging.warning("Your computer is using systemd-resolved for DNS. Thyra needs dnsmasq to redirect .massa domains to localhost.")
            user_answer = input("Do you want to install dnsmasq and configure it to redirect .massa domains to localhost? [y/n] ")
            user_answer = user_answer.lower()
            if user_answer == "y" or user_answer == "yes":
                self.executeCommand("sudo apt install -y dnsmasq", True)
                self.configureNetworkManager()
                self.configureDNSMasq()
            else:
                self.printErrorAndExit("Aborting installation.")
        else:
            logging.warning(f"Unsupported DNS application: {runningDNS}")

    def generateCACertificate(self):
        stdout, _stderr = self.executeCommand("firefox -v && echo \"firefox installed\"", True, allow_failure=True)
        if stdout and "installed" in stdout:
            self.executeCommand("sudo apt-get install -y libnss3-tools", True)

        super().generateCACertificate()
    
    def _moveFile(self, file, destination):
        self.executeCommand(f"sudo mv {file} {destination}", True)

    def _deleteFile(self, file):
        self.executeCommand(f"sudo rm {file}", True)

    """
    Override installThyraServer defined at installer level.
    """
    def installThyraServer(self):
        super().installThyraServer()
        self.executeCommand(f"sudo setcap CAP_NET_BIND_SERVICE=+eip {self.MASSASTATION_INSTALL_FOLDER_PATH}/{self.MASSASTATION_SERVER_FILENAME}", True)

if __name__ == "__main__":
    if platform.system() != "Linux":
        LinuxInstaller().printErrorAndExit("This script is only compatible with Linux.")

    LinuxInstaller().startInstall()
