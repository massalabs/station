import logging
import platform
import shutil
import ssl
import os
import subprocess

import urllib.request
from urllib.error import URLError

class Installer:
    THYRA_INSTALL_FOLDER_PATH = os.path.expanduser("~")
    THYRA_CONFIG_FOLDER_PATH = os.path.join(os.path.expanduser("~"), ".config", "thyra")
    THYRA_PLUGINS_PATH = os.path.join(THYRA_CONFIG_FOLDER_PATH, "plugins" )

    THYRA_SERVER_URL = ""
    THYRA_SERVER_FILENAME = ""

    THYRA_APP_URL = ""
    THYRA_APP_FILENAME = ""

    MKCERT_URL = ""
    MKCERT_FILENAME = ""
    CERTIFICATIONS_FOLDER_PATH = os.path.join(THYRA_CONFIG_FOLDER_PATH, "certs")

    def __init__(self):
        logging.basicConfig(filename='thyra_installer.log', level=logging.INFO)

        console = logging.StreamHandler()
        console.setLevel(logging.INFO)

        formatter = logging.Formatter('%(name)-12s: %(levelname)-8s %(message)s')
        console.setFormatter(formatter)

        logging.getLogger('').addHandler(console)
        pass

    def printErrorAndExit(self, error):
        logging.error(error)
        if platform.system() == "Windows":
            os.system("pause")
        os._exit(-1)

    def executeCommand(self, command, shell=False) -> str:
        logging.debug(f'Executing command: {command}')
        process = subprocess.Popen(command, shell=shell, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        stdout, stderr = process.communicate()

        if process.returncode != 0:
            self.printErrorAndExit(f"Command failed with error code {process.returncode}: {stderr.decode('utf-8')}")

        return stdout.decode("utf-8")

    def downloadFile(self, url, filename):
        logging.debug(f'Downloading {filename} from {url}')
        try:
            sslContext = None
            if platform.system() != "Windows":
                import certifi
                sslContext = ssl.create_default_context(cafile=certifi.where()) 
            with urllib.request.urlopen(url, context=sslContext) as response, open(filename, 'wb') as out_file:
                shutil.copyfileobj(response, out_file)
        except URLError as err:
            logging.info("Failed to download " + filename + " :")
            self.printErrorAndExit(err)
        logging.info(filename + " downloaded successfully")

    """
    This method must be implemented by the installer of each platform to install and configure the DNS server.
    """
    def setupDNS(self):
        pass

    def generateCertificate(self):
        if not os.path.exists(self.CERTIFICATIONS_FOLDER_PATH):
            try:
                os.makedirs(self.CERTIFICATIONS_FOLDER_PATH)
            except OSError as err:
                self.printErrorAndExit(f"Error while creating certs folder: {err}")

        self.downloadFile(self.MKCERT_URL, self.MKCERT_FILENAME)
        os.chmod(self.MKCERT_FILENAME, 0o755)

        self.executeCommand([
            os.path.join(os.getcwd(), self.MKCERT_FILENAME), 
            "--install"])
        self.executeCommand([
            os.path.join(os.getcwd() , self.MKCERT_FILENAME),
            "--cert-file", os.path.join(self.CERTIFICATIONS_FOLDER_PATH , "cert.pem"),
            "--key-file", os.path.join(self.CERTIFICATIONS_FOLDER_PATH, "cert-key.pem"),
            "my.massa"])

        try:
            os.remove(self.MKCERT_FILENAME)
        except OSError as err:
            self.printErrorAndExit(f"Error while deleting mkcert binary: {err}")
        logging.info("HTTPS certificate successfully generated")

    def installThyraServer(self):
        self.downloadFile(self.THYRA_SERVER_URL, self.THYRA_SERVER_FILENAME)
        os.chmod(self.THYRA_SERVER_FILENAME, 0o755)
        if os.getcwd() != self.THYRA_INSTALL_FOLDER_PATH:
            try:
                thyra_server_path = os.path.join(self.THYRA_INSTALL_FOLDER_PATH, self.THYRA_SERVER_FILENAME)
                if os.path.exists(thyra_server_path):
                    os.remove(thyra_server_path)
                shutil.move(self.THYRA_SERVER_FILENAME, self.THYRA_INSTALL_FOLDER_PATH)
            except OSError as err:
                self.printErrorAndExit(f"Error while moving thyra server binary: {err}")
        logging.info("Thyra server installed successfully")

    def installThyraApp(self):
        self.downloadFile(self.THYRA_APP_URL, self.THYRA_APP_FILENAME)
        os.chmod(self.THYRA_APP_FILENAME, 0o755)
        if os.getcwd() != self.THYRA_INSTALL_FOLDER_PATH:
            try:
                thyra_app_path = os.path.join(self.THYRA_INSTALL_FOLDER_PATH, self.THYRA_APP_FILENAME)
                if os.path.exists(thyra_app_path):
                    os.remove(thyra_app_path)
                shutil.move(self.THYRA_APP_FILENAME, self.THYRA_INSTALL_FOLDER_PATH)
            except OSError as err:
                self.printErrorAndExit(f"Error while moving thyra app binary: {err}")
        logging.info("Thyra app installed successfully")

    def startInstall(self):
        logging.info("Starting installation of thyra")
        self.installThyraServer()
        self.installThyraApp()
        self.setupDNS()
        self.generateCertificate()
        logging.info("Thyra installed successfully")
