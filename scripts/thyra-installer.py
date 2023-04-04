
### legacy: to be deleted

import platform
import ssl
import urllib.request
import subprocess
import shutil
import os
import zipfile
import ctypes
import logging
from urllib.error import URLError
from tarfile import ReadError

import logging
import os

# This file is to be bundled with pyinstaller in order to produce a .exe that can run on Windows without Python installed.

# set up logging to file
logging.basicConfig(level=logging.INFO,
                    filename=os.path.join(os.path.expanduser("~"), "thyra_installer.log"),
                    filemode='w')

# define a Handler which writes INFO messages or higher to the sys.stderr
console = logging.StreamHandler()
console.setLevel(logging.INFO)

# set a format which is simpler for console use
formatter = logging.Formatter('%(name)-12s: %(levelname)-8s %(message)s')
console.setFormatter(formatter)

# add the handler to the root logger
logging.getLogger('').addHandler(console)

logging.info('Starting thyra installer . . .')

# General
THYRA_SERVER_URL = ""
THYRA_SERVER_FILENAME = ""
THYRA_APP_URL = ""
THYRA_APP_FILENAME = ""
THYRA_CONFIG_FOLDER_PATH = os.path.join(os.path.expanduser("~"), ".config", "thyra")
THYRA_PLUGINS_PATH = os.path.join(THYRA_CONFIG_FOLDER_PATH, "plugins" )

USER_HOME_FOLDER = os.path.expanduser("~")

# Windows // Acrylic
ACRYLIC_DNS_PROXY_URL = "https://sourceforge.net/projects/acrylic/files/Acrylic/2.1.1/Acrylic-Portable.zip/download"
ACRYLIC_DNS_PROXY_FILENAME = "Acrylic-Portable.zip"
ACRYLIC_HOST_FILE = "AcrylicHosts.txt"
DEFAULT_ACRYLIC_PATH = "C:\Program Files (x86)\Acrylic DNS Proxy"

# Certifications
MKCERT_URL = "https://dl.filippo.io/mkcert/latest?for=windows/amd64"
MKCERT_FILENAME = "mkcert.exe"
CERTIFICATIONS_FOLDER = os.path.join(THYRA_CONFIG_FOLDER_PATH, "certs")

# Global variables
def setThyraGlobals():
    global THYRA_SERVER_URL, THYRA_SERVER_FILENAME

    if platform.system() == "Windows":
        THYRA_SERVER_FILENAME = "thyra-server.exe"
        THYRA_SERVER_URL = "https://github.com/massalabs/thyra/releases/latest/download/thyra-server_windows_amd64"
    elif platform.system() == "Darwin":
        THYRA_SERVER_FILENAME = "thyra-server"
        if platform.machine() == "arm64":
            THYRA_SERVER_URL = "https://github.com/massalabs/thyra/releases/latest/download/thyra-server_darwin_arm64"
        elif platform.machine() == "x86_64":
            THYRA_SERVER_URL = "https://github.com/massalabs/thyra/releases/latest/download/thyra-server_darwin_amd64"
    else:
        printErrorAndExit("Unsupported platform: " + platform.system())

def setThyraAppGlobals():
    global THYRA_APP_URL, THYRA_APP_FILENAME

    if platform.system() == "Windows":
            THYRA_APP_FILENAME = "thyra-app.exe"
            THYRA_APP_URL = "https://github.com/massalabs/Thyra-Menu-Bar-App/releases/latest/download/ThyraApp_windows-amd64.exe"
    elif platform.system() == "Darwin":
            THYRA_APP_FILENAME = "thyra-app"
            if platform.machine() == "arm64":
                THYRA_APP_URL = "https://github.com/massalabs/Thyra-Menu-Bar-App/releases/latest/download/ThyraApp_darwin-arm64"
            elif platform.machine() == "x86_64":
                THYRA_APP_URL = "https://github.com/massalabs/Thyra-Menu-Bar-App/releases/latest/download/ThyraApp_darwin-amd64"
    else:
        printErrorAndExit("Unsupported platform: " + platform.system())

def setMKCertGlobals():
    global MKCERT_URL, MKCERT_FILENAME

    if platform.system() == "Windows":
            MKCERT_URL = "https://dl.filippo.io/mkcert/latest?for=windows/amd64"
            MKCERT_FILENAME = "mkcert.exe"
    elif platform.system() == "Darwin":
            MKCERT_FILENAME = "mkcert"
            if platform.machine() == "arm64":
                MKCERT_URL = "https://dl.filippo.io/mkcert/latest?for=darwin/arm64"
            elif platform.machine() == "x86_64":
                MKCERT_URL = "https://dl.filippo.io/mkcert/latest?for=darwin/amd64"
    else:
        printErrorAndExit("Unsupported platform: " + platform.system())

# Windows
def unzipAcrylic():
    try: 
        os.mkdir(DEFAULT_ACRYLIC_PATH)
        shutil.move(ACRYLIC_DNS_PROXY_FILENAME, DEFAULT_ACRYLIC_PATH)
        os.chdir(DEFAULT_ACRYLIC_PATH)
    except OSError as err: 
        printErrorAndExit(err)
    try:
        with zipfile.ZipFile(DEFAULT_ACRYLIC_PATH + "\\" +  ACRYLIC_DNS_PROXY_FILENAME, 'r') as zip_ref:
            zip_ref.extractall(DEFAULT_ACRYLIC_PATH)
        logging.info("Acrylic unzipped")
    except ReadError as err:
        printErrorAndExit(err)

def setupDNS():
    commandOutput = executeOSCommandOrFile("wmic nic where \"netenabled=true\" get netconnectionID", True)

    networkAdpatersNames = list(filter(None, commandOutput.split('\n')))
    networkAdpatersNames.pop(0)
    networkAdpatersNames = [name.strip() for name in networkAdpatersNames]

    for name in networkAdpatersNames:
        executeOSCommandOrFile("NETSH interface ipv4 set dnsservers " + name + " static 127.0.0.1 primary", True)
        executeOSCommandOrFile("NETSH interface ipv6 set dnsservers " + name + " static ::1 primary", True)
    try:
        os.remove(DEFAULT_ACRYLIC_PATH + "\\" +  ACRYLIC_DNS_PROXY_FILENAME)
    except OSError as err:
        printErrorAndExit(err)

def configureAcrylic():
    logging.info("Configuring Acrylic...")
    f = open(DEFAULT_ACRYLIC_PATH + "\\" + ACRYLIC_HOST_FILE, "r+")
    if f.read().find("127.0.0.1 *.massa") != -1:
        f.close()
        return
    f.write("\n127.0.0.1 *.massa")
    f.close()

    executeOSCommandOrFile("NET STOP AcrylicDNSProxySvc", True)
    executeOSCommandOrFile("NET START AcrylicDNSProxySvc", True)

def generateCertificate():
    if os.path.isdir(THYRA_CONFIG_FOLDER_PATH) is False:
        try:
            os.mkdir(THYRA_CONFIG_FOLDER_PATH)
        except OSError as err:
            printErrorAndExit("Error while creating config folder: " + err)

    if os.path.isdir(CERTIFICATIONS_FOLDER) is False:
        try:
            os.mkdir(CERTIFICATIONS_FOLDER)
        except OSError as err:
            printErrorAndExit("Error while creating certificates folder: " + err)

    downloadFile(MKCERT_URL, MKCERT_FILENAME)
    os.chmod(MKCERT_FILENAME, 755)

    executeOSCommandOrFile([os.path.join(os.getcwd(), MKCERT_FILENAME), "--install"], False)
    executeOSCommandOrFile([
        os.path.join(os.getcwd() , MKCERT_FILENAME),
        "--cert-file", os.path.join(CERTIFICATIONS_FOLDER , "cert.pem"),
        "--key-file", os.path.join(CERTIFICATIONS_FOLDER, "cert-key.pem"),
        "my.massa"], False)
    try:
        os.remove(MKCERT_FILENAME)
    except OSError as err:
        printErrorAndExit(err)
    logging.info("MKcert and HTTPS certificates successfully setup")

def isAdmin():
    return ctypes.windll.shell32.IsUserAnAdmin() != 0

# MacOS
def configureDNSMasq():
    print("Configuring DNSMasq...")
    
    executeOSCommandOrFile("sudo bash -c 'echo ""address=/.massa/127.0.0.1"" > $(brew --prefix)/etc/dnsmasq.d/massa.conf'", True, True)
    executeOSCommandOrFile("sudo mkdir -p /etc/resolver", True, True)
    executeOSCommandOrFile("sudo bash -c 'echo ""nameserver 127.0.0.1"" > /etc/resolver/massa'", True, True)

    print("Restarting DNSMasq...")
    executeOSCommandOrFile("sudo brew services restart dnsmasq", True, True)


# Generic
def downloadFile(url, filename):
    logging.info("Downloading " + filename + "...")
    try:
        sslContext = None
        if platform.system() != "Windows":
            import certifi
            sslContext = ssl.create_default_context(cafile=certifi.where()) 
        with urllib.request.urlopen(url, context=sslContext) as response, open(filename, 'wb') as out_file:
            shutil.copyfileobj(response, out_file)
    except URLError as err:
        logging.info("Failed to download " + filename + " :")
        printErrorAndExit(err)
    logging.info(filename + " downloaded successfully")

def executeOSCommandOrFile(command, decodeBinary, shell=False):
    process = subprocess.Popen(command,
                            stdout=subprocess.PIPE,
                            stderr=subprocess.PIPE,
                            universal_newlines=decodeBinary,
                            shell=shell)
    stdout, stderr = process.communicate()

    if stderr != None and stderr != "" and process.returncode != 0:
        printErrorAndExit(f"Error encountered while executing : {command} :\n{stderr}")
    return stdout

def printErrorAndExit(error):
    print(error)
    if platform.system() == "Windows":
        os.system("pause")
    os._exit(-1)

def main():
    if not os.path.exists(THYRA_CONFIG_FOLDER_PATH):
        os.makedirs(THYRA_CONFIG_FOLDER_PATH)
    if not os.path.exists(THYRA_PLUGINS_PATH):
        os.makedirs(THYRA_PLUGINS_PATH)    
    if platform.system() == "Windows" and isAdmin() == False:
        printErrorAndExit("Couldn't detect admin rights. Please execute this script as an administator.")

    setThyraGlobals()
    setThyraAppGlobals()
    setMKCertGlobals()

    downloadFile(THYRA_SERVER_URL, THYRA_SERVER_FILENAME)
    downloadFile(THYRA_APP_URL, THYRA_APP_FILENAME)
    try:
        os.chmod(THYRA_SERVER_FILENAME, 0o755)
        os.chmod(THYRA_APP_FILENAME, 0o755)

        if os.getcwd() != USER_HOME_FOLDER:
            thyra_home_path = os.path.join(USER_HOME_FOLDER, THYRA_SERVER_FILENAME)
            thyra_app_home_path = os.path.join(USER_HOME_FOLDER, THYRA_APP_FILENAME)
            if os.path.exists(thyra_home_path):
                os.remove(thyra_home_path)
            if os.path.exists(thyra_app_home_path):
                os.remove(thyra_app_home_path)

            shutil.move(THYRA_SERVER_FILENAME, USER_HOME_FOLDER)
            shutil.move(THYRA_APP_FILENAME, USER_HOME_FOLDER)
    except OSError as err:
        logging.error("Error while installing Thyra:")
        printErrorAndExit(err)

    if platform.system() == "Windows":
        if os.path.exists(DEFAULT_ACRYLIC_PATH):
            logging.info("Acrylic DNS Proxy is already installed")
        else:
            downloadFile(ACRYLIC_DNS_PROXY_URL, ACRYLIC_DNS_PROXY_FILENAME)
            unzipAcrylic()
            executeOSCommandOrFile(DEFAULT_ACRYLIC_PATH + "\InstallAcrylicService.bat", True)
            setupDNS()
        configureAcrylic()
    elif platform.system() == "Darwin":
        runningDNS = executeOSCommandOrFile("sudo lsof -i :53 | sed -n 2p | sed 's/[[:space:]].*$//'", True, True)
        runningDNS = runningDNS[:-1]
        if runningDNS == "dnsmasq":
            logging.info("dnsmasq is already installed")
            configureDNSMasq()
        elif runningDNS == "":
            logging.info("Installing dnsmasq...")
            executeOSCommandOrFile("brew install dnsmasq", False)
            configureDNSMasq()
        else:
            logging.info(runningDNS)
            printErrorAndExit("Unsupported DNS server")
    else:
        printErrorAndExit("Unsupported platform: " + platform.system())

    generateCertificate()

    logging.info("Thyra has been successfully installed! Executable is located at : " + USER_HOME_FOLDER)

    THYRA_APP_PATH = os.path.join(USER_HOME_FOLDER, THYRA_APP_FILENAME)
    if os.path.exists(THYRA_APP_PATH):
        subprocess.Popen([THYRA_APP_PATH], start_new_session=True)
        logging.info("You can start using thyra from the menu bar located on the bottom of your screen")

    if platform.system() == "Windows":
        os.system("pause")
    os._exit(0)

main()
