import urllib.request
import subprocess
import shutil
import os
import zipfile
import ctypes
from urllib.error import URLError
from tarfile import ReadError

# This file is to be bundled with pyinstaller in order to produce a .exe that can run on Windows without Python installed.

THYRA_URL = "https://github.com/massalabs/thyra/releases/latest/download/thyra-server_windows_amd64"
THYRA_APP = "https://github.com/massalabs/Thyra-Menu-Bar-App/releases/download/v0.0.1/ThyraApp_windows-amd64.exe"
THYRA_APP_FILENAME = "ThyraApp_windows-amd64.exe"
THYRA_FILENAME = "thyra-server.exe"
THYRA_CONFIG_FOLDER_PATH = os.path.expanduser("~") + "\\.config\\thyra"
ACRYLIC_DNS_PROXY_URL = "https://sourceforge.net/projects/acrylic/files/Acrylic/2.1.1/Acrylic-Portable.zip/download"
ACRYLIC_DNS_PROXY_FILENAME = "Acrylic-Portable.zip"
ACRYLIC_HOST_FILE = "AcrylicHosts.txt"
DEFAULT_ACRYLIC_PATH = "C:\Program Files (x86)\Acrylic DNS Proxy"
MKCERT_URL = "https://dl.filippo.io/mkcert/latest?for=windows/amd64"
MKCERT_FILENAME = "mkcert.exe"
CERTIFICATIONS_FOLDER = THYRA_CONFIG_FOLDER_PATH + "\\certs"
USER_HOME_FOLDER = os.path.expanduser("~")

def downloadFile(url, filename):
    print("Downloading " + filename + "...")
    try:
        with urllib.request.urlopen(url) as response, open(filename, 'wb') as out_file:
            shutil.copyfileobj(response, out_file)
    except URLError as err:
        print("Failed to download " + filename + " :")
        printErrorAndExit(err)
    print(filename + " downloaded successfully")

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
        print("Acrylic unzipped")
    except ReadError as err:
        printErrorAndExit(err)

def executeOSCommandOrFile(command, decodeBinary):
    process = subprocess.Popen(command,
                            stdout=subprocess.PIPE,
                            stderr=subprocess.PIPE,
                            universal_newlines=decodeBinary)
    stdout, stderr = process.communicate()

    if stderr != None and stderr != "" and process.returncode != 0:
        printErrorAndExit("Error encountered while executing : " + str(command) + " :\n" + stderr)
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
    print("Configuring Acrylic...")
    f = open(DEFAULT_ACRYLIC_PATH + "\\" + ACRYLIC_HOST_FILE, "r+")
    if f.read().find("127.0.0.1 *.massa") != -1:
        f.close()
        return
    f.write("\n127.0.0.1 *.massa")
    f.close()

    executeOSCommandOrFile("NET STOP AcrylicDNSProxySvc", True)
    executeOSCommandOrFile("NET START AcrylicDNSProxySvc", True)

def setupMkCerts():
    if os.path.isdir(THYRA_CONFIG_FOLDER_PATH) is False:
        try:
            os.mkdir(THYRA_CONFIG_FOLDER_PATH)
            os.mkdir(THYRA_CONFIG_FOLDER_PATH + "\\certs")
        except OSError as err:
            printErrorAndExit(err)

    downloadFile(MKCERT_URL, MKCERT_FILENAME)
    executeOSCommandOrFile([os.getcwd() + "\\" + MKCERT_FILENAME, "--install"], False)
    executeOSCommandOrFile([
        os.getcwd() + "\\" + MKCERT_FILENAME,
        "--cert-file",
        CERTIFICATIONS_FOLDER + "\\cert.pem",
        "--key-file",
        CERTIFICATIONS_FOLDER  + "\\cert-key.pem",
        "my.massa"], False)
    try:
        os.remove(MKCERT_FILENAME)
    except OSError as err:
        printErrorAndExit(err)
    print("MKcert and HTTPS certificates successfully setup")

def isAdmin():
    return ctypes.windll.shell32.IsUserAnAdmin() != 0

def printErrorAndExit(error):
    print(error)
    os.system("pause")
    os._exit(-1)

def main():
    if not os.path.exists(THYRA_CONFIG_FOLDER_PATH):
        os.makedirs(THYRA_CONFIG_FOLDER_PATH)
    if isAdmin() == False:
        printErrorAndExit("Couldn't detect admin rights. Please execute this script as an administator.")
    downloadFile(THYRA_URL, THYRA_FILENAME)
    downloadFile(THYRA_APP, THYRA_APP_FILENAME)
    try:
        shutil.move(THYRA_FILENAME, USER_HOME_FOLDER)
        shutil.move(THYRA_APP_FILENAME, USER_HOME_FOLDER)
    except OSError as err:
        printErrorAndExit(err)    
    if os.path.exists(DEFAULT_ACRYLIC_PATH):
        print("Acrylic DNS Proxy is already installed")
    else:
        downloadFile(ACRYLIC_DNS_PROXY_URL, ACRYLIC_DNS_PROXY_FILENAME)
        unzipAcrylic()
        executeOSCommandOrFile(DEFAULT_ACRYLIC_PATH + "\InstallAcrylicService.bat", True)
        setupDNS()
    configureAcrylic()
    setupMkCerts()
    THYRA_APP_PATH = os.path.join(USER_HOME_FOLDER, THYRA_APP_FILENAME)
    print("Thyra has been successfully installed! Executable is located at : " + USER_HOME_FOLDER)
    if os.path.exists(THYRA_APP_PATH):
        executeOSCommandOrFile(THYRA_APP_PATH, True)
        print("You can start using thyra from the menu bar located on the bottom of your screen")
    os.system("pause")
    os._exit(0)

main()