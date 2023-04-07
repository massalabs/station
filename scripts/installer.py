import locale
import logging
import os
import platform
import shutil
import socket
import ssl
import subprocess
import traceback
import urllib.request
from urllib.error import URLError

logging.basicConfig(
    level=os.environ.get("LOGLEVEL", "INFO"),
    format='%(name)-12s: %(levelname)-8s %(message)s',
    handlers=[
        logging.FileHandler("thyra_installer.log", mode='w'),
        logging.StreamHandler()
    ]
)

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
    CERTIFICATION_FILENAME = "cert.pem"
    CERTIFICATION_KEY_FILENAME = "cert-key.pem"

    def __init__(self):
        pass

    """
    Prints the error given in parameter and exits the program
    """
    @staticmethod
    def printErrorAndExit(error):
        logging.error(error)
        os._exit(-1)

    """
    Executes the command given in parameter and returns the output of the command
    """
    def executeCommand(self, command, shell=False, allow_failure=False) -> tuple[str, str]:
        logging.debug(f'Executing command: {command}')
        try:
            process = subprocess.Popen(command, shell=shell,
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE
            )
            stdout, stderr = process.communicate()
            try:

                if stdout is not None and len(stdout) > 0:
                    logging.debug(f'Command output: {stdout.decode(locale.getpreferredencoding(False))}')

                if allow_failure == False and process.returncode != 0:
                    self.printErrorAndExit(f"Command failed with error : {stderr.decode(locale.getpreferredencoding(False))}")

                return (stdout.decode(locale.getpreferredencoding(False)), stderr.decode(locale.getpreferredencoding(False)))
            except :
                return(stdout, stderr)
        except OSError as err:
            self.printErrorAndExit(f"Error while executing command: {err}")

    """
    Downloads the file at the given url and saves it to the given filename
    """
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
    Checks if my.massa is resolved. If it is, the DNS server is well installed and configured.
    """
    def shouldInstallDNS(self) -> bool:
        try:
            socket.gethostbyname("my.massa")
            return False
        except:
            return True

    """
    This method must be implemented by the installer of each platform to install and configure the DNS server.
    """
    def setupDNS(self):
        pass

    """
    Checks if the HTTPS certificate for my.massa are already present in the thyra config folder.
    The filenames of the certificate and the key are defined in the constants CERTIFICATION_FILENAME and CERTIFICATION_KEY_FILENAME.
    """
    def isCACertificateInstalled(self) -> bool:
        return os.path.exists(os.path.join(self.CERTIFICATIONS_FOLDER_PATH, self.CERTIFICATION_FILENAME)) and os.path.exists(os.path.join(self.CERTIFICATIONS_FOLDER_PATH, self.CERTIFICATION_KEY_FILENAME))

    """
    Installs a local Certificate Authority and generates a HTTPS certificate for my.massa.
    """
    def generateCACertificate(self):
        logging.info("Generating CA certificate and adding it to the browsers' CA list")

        if not os.path.exists(self.CERTIFICATIONS_FOLDER_PATH):
            try:
                os.makedirs(self.CERTIFICATIONS_FOLDER_PATH)
            except OSError as err:
                self.printErrorAndExit(f"Error while creating certs folder: {err}")

        self.downloadFile(self.MKCERT_URL, self.MKCERT_FILENAME)
        os.chmod(self.MKCERT_FILENAME, 0o755)

        stdout, stderr = self.executeCommand([
            os.path.join(os.getcwd(), self.MKCERT_FILENAME), 
            "--install"])
        if stderr is not None and len(stderr) > 0:
            logging.info(stderr)

        stdout, stderr = self.executeCommand([
            os.path.join(os.getcwd() , self.MKCERT_FILENAME),
            "--cert-file", os.path.join(self.CERTIFICATIONS_FOLDER_PATH, self.CERTIFICATION_FILENAME),
            "--key-file", os.path.join(self.CERTIFICATIONS_FOLDER_PATH, self.CERTIFICATION_KEY_FILENAME),
            "my.massa"])
        if stderr is not None and len(stderr) > 0:
            logging.info(stderr)

        try:
            os.remove(self.MKCERT_FILENAME)
        except OSError as err:
            self.printErrorAndExit(f"Error while deleting mkcert binary: {err}")
        logging.info("CA certificate successfully generated")

    def _moveFile(self, file, destination):
        shutil.move(file, destination)
    
    def _deleteFile(self, file):
        os.remove(file)

    """
    Downloads and installs a binary from the given url and stores it in the given install path.
    """
    def installBinary(self, install_path, binary_url, binary_filename):
        logging.debug(f"Installing {binary_filename} from {binary_url}")
        self.downloadFile(binary_url, binary_filename)
        os.chmod(binary_filename, 0o755)
        if os.getcwd() != install_path:
            try:
                thyra_server_path = os.path.join(install_path, binary_filename)
                if os.path.exists(thyra_server_path):
                    self._deleteFile(thyra_server_path)
                self._moveFile(binary_filename, install_path)
            except OSError as err:
                self.printErrorAndExit(f"Error while moving {binary_filename} binary: {err}")
        logging.debug(f"{binary_filename} successfully installed")

    """
    Downloads thyra server and stores it in the thyra install folder.
    Can be overriden by the installer of a specific platform if needed.
    """
    def installThyraServer(self):
        self.installBinary(self.THYRA_INSTALL_FOLDER_PATH, self.THYRA_SERVER_URL, self.THYRA_SERVER_FILENAME)
        logging.info("Thyra server installed successfully")

    """
    Downloads thyra app and stores it in the thyra install folder.
    Can be overriden by the installer of a specific platform if needed.
    """
    def installThyraApp(self):
        self.installBinary(self.THYRA_INSTALL_FOLDER_PATH, self.THYRA_APP_URL, self.THYRA_APP_FILENAME)
        logging.info("Thyra app installed successfully")

    def createConfigFolder(self):
        logging.info("Creating config folder")
        if not os.path.exists(self.THYRA_CONFIG_FOLDER_PATH):
            try:
                os.makedirs(self.THYRA_CONFIG_FOLDER_PATH)
                os.makedirs(self.THYRA_PLUGINS_PATH)
            except OSError as err:
                self.printErrorAndExit(f"Error while creating config folder: {err}")

    def startThyraApp(self):
        thyra_app_path = os.path.join(self.THYRA_INSTALL_FOLDER_PATH, self.THYRA_APP_FILENAME)
        if os.path.exists(thyra_app_path):
            subprocess.Popen([thyra_app_path], start_new_session=True)
            logging.info("Thyra App will now start. You can right click on the tray icon to start Thyra.")

    """
    Installs thyra server, thyra app and a DNS server.
    """
    def startInstall(self):
        logging.info("Starting installation of thyra")
        try:
            self.createConfigFolder()
            self.installThyraServer()
            self.installThyraApp()
            if self.shouldInstallDNS():
                self.setupDNS()
            else:
                logging.info("DNS server already installed.")
            if not self.isCACertificateInstalled():
                self.generateCACertificate()
            else:
                logging.info("CA certificate already installed.")

            logging.info("Thyra installed successfully !")
            logging.info(f"{self.THYRA_SERVER_FILENAME} and {self.THYRA_APP_FILENAME} has been installed in {self.THYRA_INSTALL_FOLDER_PATH}")

            self.startThyraApp()
        except Exception as e:
            logging.error(traceback.format_exc())
        logging.info("You can now close this window.")
