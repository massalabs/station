import platform
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

# This file is to be bundled with pyinstaller in order to produce a .exe that can run on Windows without Python installed.
THYRA_URL = "https://github.com/massalabs/thyra/releases/latest/download/thyra-server_windows_amd64"
THYRA_FILENAME = "thyra-server.exe"
THYRA_APP = "https://github.com/massalabs/Thyra-Menu-Bar-App/releases/latest/download/ThyraApp_windows-amd64.exe"
THYRA_APP_FILENAME = "ThyraApp_windows-amd64.exe"

THYRA_URL = ""
THYRA_FILENAME = ""
THYRA_CONFIG_FOLDER_PATH = os.path.join(os.path.expanduser("~"), ".config", "thyra")

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


def setThyraGlobals():
    global THYRA_URL, THYRA_FILENAME

    match platform.system():
        case "Windows":
            THYRA_FILENAME = "thyra-server.exe"
            THYRA_URL = "https://github.com/massalabs/thyra/releases/latest/download/thyra-server_windows_amd64"
        case "Darwin":
            THYRA_FILENAME = "thyra-server"
            match platform.machine():
                case "arm64":
                    THYRA_URL = "https://github.com/massalabs/thyra/releases/latest/download/thyra-server_darwin_arm64"
                case "x86_64":
                    THYRA_URL = "https://github.com/massalabs/thyra/releases/latest/download/thyra-server_darwin_amd64"
        case _:
            printErrorAndExit("Unsupported platform: " + platform.system())

def setMKCertGlobals():
    global MKCERT_URL, MKCERT_FILENAME

    match platform.system():
        case "Windows":
            MKCERT_URL = "https://dl.filippo.io/mkcert/latest?for=windows/amd64"
            MKCERT_FILENAME = "mkcert.exe"
        case "Darwin":
            MKCERT_FILENAME = "mkcert"
            match platform.machine():
                case "arm64":
                    MKCERT_URL = "https://dl.filippo.io/mkcert/latest?for=darwin/arm64"
                case "x86_64":
                    MKCERT_URL = "https://dl.filippo.io/mkcert/latest?for=darwin/amd64"
        case _:
            printErrorAndExit("Unsupported platform: " + platform.system())

def downloadFile(url, filename):
    logging.info("Downloading " + filename + "...")
    try:
        with urllib.request.urlopen(url) as response, open(filename, 'wb') as out_file:
            shutil.copyfileobj(response, out_file)
    except URLError as err:
        logging.info("Failed to download " + filename + " :")
        printErrorAndExit(err)
    logging.info(filename + " downloaded successfully")

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

def executeOSCommandOrFile(command, decodeBinary):
    process = subprocess.Popen(command,
                            stdout=subprocess.PIPE,
                            stderr=subprocess.PIPE,
                            universal_newlines=decodeBinary)
    stdout, stderr = process.communicate()

    if stderr != None and stderr != "" and process.returncode != 0:
        printErrorAndExit(f"Error encountered while executing : {command} :\n{stderr}")
    return stdout

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

def printErrorAndExit(error):
    print(error)
    if platform.system() == "Windows":
        os.system("pause")
    os._exit(-1)

def main():
    if not os.path.exists(THYRA_CONFIG_FOLDER_PATH):
        os.makedirs(THYRA_CONFIG_FOLDER_PATH)
    if platform.system() == "Windows" and isAdmin() == False:
        printErrorAndExit("Couldn't detect admin rights. Please execute this script as an administator.")

    setThyraGlobals()
    setMKCertGlobals()

    # DEBUG
    print(THYRA_FILENAME, THYRA_URL)

    downloadFile(THYRA_URL, THYRA_FILENAME)
    downloadFile(THYRA_APP, THYRA_APP_FILENAME)
    try:
        os.chmod(THYRA_FILENAME, 755)

        thyra_home_path = os.path.join(USER_HOME_FOLDER, THYRA_FILENAME)
        if os.path.exists(thyra_home_path):
            os.remove(thyra_home_path)

        shutil.move(THYRA_FILENAME, USER_HOME_FOLDER)
        shutil.move(THYRA_APP_FILENAME, USER_HOME_FOLDER)
    except OSError as err:
        printErrorAndExit(err)

    match platform.system():
        case "Windows":
            if os.path.exists(DEFAULT_ACRYLIC_PATH):
                logging.info("Acrylic DNS Proxy is already installed")
            else:
                downloadFile(ACRYLIC_DNS_PROXY_URL, ACRYLIC_DNS_PROXY_FILENAME)
                unzipAcrylic()
                executeOSCommandOrFile(DEFAULT_ACRYLIC_PATH + "\InstallAcrylicService.bat", True)
                setupDNS()
            configureAcrylic()
        case "Darwin":
            logging.info("MacOS is not supported yet")
        case _:
            printErrorAndExit("Unsupported platform: " + platform.system())

    generateCertificate()

    logging.info("Thyra has been successfully installed! Executable is located at : " + USER_HOME_FOLDER)

    THYRA_APP_PATH = os.path.join(USER_HOME_FOLDER, THYRA_APP_FILENAME)
    if os.path.exists(THYRA_APP_PATH):
        executeOSCommandOrFile(THYRA_APP_PATH, True)
        logging.info("You can start using thyra from the menu bar located on the bottom of your screen")

    if platform.system() == "Windows":
        os.system("pause")
    os._exit(0)

main()