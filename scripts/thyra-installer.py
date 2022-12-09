import urllib.request
import subprocess
from urllib.error import URLError
from tarfile import ReadError
import shutil
import os
import tarfile
import zipfile


THYRA_URL="https://github.com/massalabs/thyra/releases/latest/download/thyra-server_windows_amd64"
THYRA_FILENAME="thyra-server.exe"
ACRYLIC_DNS_PROXY_URL="https://sourceforge.net/projects/acrylic/files/Acrylic/2.1.1/Acrylic-Portable.zip/download"
ACRYLIC_DNS_PROXY_FILENAME="Acrylic-Portable.zip"
DEFAULT_ACRYLIC_PATH="C:\Program Files (x86)\Acrylic DNS Proxy"

def downloadThyra():
    print("Downloading Thyra...")
    try:
        with urllib.request.urlopen(THYRA_URL) as response, open(THYRA_FILENAME, 'wb') as out_file:
            shutil.copyfileobj(response, out_file)
    except URLError as e:
        print("Failed to download Thyra")
        exit(-1)
    print("Thyra downloaded successfully")

def downloadAcrylic():
    print("Downloading Acrylic DNS Proxy...")
    try:
        with urllib.request.urlopen(ACRYLIC_DNS_PROXY_URL) as response, open(ACRYLIC_DNS_PROXY_FILENAME, 'wb') as out_file:
            shutil.copyfileobj(response, out_file)
    except URLError as e:
        print("Failed to download Acrylic DNS Proxy")
        exit(-1)
    print("Acrylic DNS Proxy downloaded successfully")

def unzipAcrylic():
    try: 
        os.mkdir(DEFAULT_ACRYLIC_PATH)
        shutil.move(ACRYLIC_DNS_PROXY_FILENAME, DEFAULT_ACRYLIC_PATH)
        os.chdir(DEFAULT_ACRYLIC_PATH)
    except OSError as error: 
        print(error)
        exit(-1)
    try:
        with zipfile.ZipFile(DEFAULT_ACRYLIC_PATH + "\\" +  ACRYLIC_DNS_PROXY_FILENAME, 'r') as zip_ref:
            zip_ref.extractall(DEFAULT_ACRYLIC_PATH)
        print("Acrylic unzipped")
    except ReadError as error:
        print(error)
        exit(-1)

def installAcrylic():
    process = subprocess.Popen(DEFAULT_ACRYLIC_PATH + "\InstallAcrylicService.bat",
                            stdout=subprocess.PIPE,
                            stderr=subprocess.PIPE,
                            universal_newlines=True)
    stdout, stderr = process.communicate()
    stdout, stderr

def main():
    #downloadThyra()
    if os.path.exists(DEFAULT_ACRYLIC_PATH):
        print("Acrylic DNS Proxy is installed")
    else:
        print("Acrylic DNS Proxy not installed")
        downloadAcrylic()
        unzipAcrylic()
        installAcrylic()

main()