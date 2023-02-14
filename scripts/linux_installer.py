import logging
import platform

from installer import Installer

class LinuxInstaller(Installer):
    def __init__(self):
        super().__init__()
        self.THYRA_SERVER_FILENAME = "thyra_server"
        self.THYRA_APP_FILENAME = "thyra-app"
        self.MKCERT_FILENAME = "mkcert"

        self.THYRA_SERVER_URL = "https://github.com/massalabs/thyra/releases/latest/download/thyra-server_linux"
        self.THYRA_APP_URL = "https://github.com/massalabs/Thyra-Menu-Bar-App/releases/latest/download/ThyraApp_linux-amd64.exe"
        self.MKCERT_URL = "https://dl.filippo.io/mkcert/latest?for=linux/amd64"

    def configureNetworkManager(self):
        logging.info("Configuring NetworkManager...")
        self.executeCommand("sudo cp /etc/NetworkManager/NetworkManager.conf /etc/NetworkManager/NetworkManager.conf.backup_thyra_install", True)
        stdout, _stderr = self.executeCommand("grep '^dns=' /etc/NetworkManager/NetworkManager.conf", True)
        dns = stdout.removeprefix("dns=").strip()
        if dns == "dnsmasq":
            logging.info("NetworkManager is already configured")
        elif dns == "":
            self.executeCommand("sudo sed -i 's/^\[main\]$/\[main\]\ndns=dnsmasq/g' /etc/NetworkManager/NetworkManager.conf", True)
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
            runningDNS = runningDNS[:-1]
        
        if runningDNS == "":
            logging.info("Installing dnsmasq...")
            self.executeCommand("sudo apt-get install -y dnsmasq", True)
            self.configureNetworkManager()
            self.configureDNSMasq()
        elif runningDNS == "dnsmasq":
            logging.info("dnsmasq is already installed")
            self.configureDNSMasq()
        elif runningDNS == "systemd-r":
            logging.warn("Your computer is using systemd-resolved for DNS. Thyra needs dnsmasq to redirect .massa domains to localhost.")
            user_answer = input("Do you want to install dnsmasq and configure it to redirect .massa domains to localhost? [y/n] ")
            user_answer = user_answer.lower()
            if user_answer == "y" or user_answer == "yes":
                self.configureNetworkManager()
                self.configureDNSMasq()
            else:
                self.printErrorAndExit("Aborting installation.")
        else:
            logging.warn(f"Unsupported DNS application: {runningDNS}")


if __name__ == "__main__":
    if platform.system() != "Linux":
        LinuxInstaller().printErrorAndExit("This script is only compatible with Linux.")

    LinuxInstaller().startInstall()
