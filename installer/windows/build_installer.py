#!/usr/bin/env python3

"""
This script generates a .msi file for the installation of MassaStation on Windows.
"""

import argparse
import os
import re
import shutil
import subprocess
import sys
import urllib.request
import zipfile

BUILD_DIR = "buildmsi"

VERSION = "0.0.0"

# Binaries to be included in the installer
SERVER_BINARIES = "massastation-server.exe"
APP_BINARIES = "massastation-app.exe"
ACRYLIC_ZIP = "acrylic.zip"

# Scripts to be included in the installer
ACRYLIC_CONFIG_SCRIPT = "configure_acrylic.bat"
NIC_CONFIG_SCRIPT = "configure_network_interfaces.bat"
GEN_CERT_SCRIPT = "generate_certificate.bat"

WIX_DIR = "wixtoolset"


def download_file(url, filename):
    """
    Download a given file from a given URL.
    """
    urllib.request.urlretrieve(url, filename)


def move_binaries():
    """
    Move the binaries and scripts to the build directory.

    If the build directory already exists, it is deleted first.
    """
    if os.path.exists(BUILD_DIR):
        shutil.rmtree(BUILD_DIR)

    os.makedirs(BUILD_DIR)

    os.rename(SERVER_BINARIES, os.path.join(BUILD_DIR, SERVER_BINARIES))
    os.rename(APP_BINARIES, os.path.join(BUILD_DIR, APP_BINARIES))
    os.rename(ACRYLIC_ZIP, os.path.join(BUILD_DIR, ACRYLIC_ZIP))

    shutil.copy(
        os.path.join("windows", "scripts", ACRYLIC_CONFIG_SCRIPT),
        os.path.join(BUILD_DIR, ACRYLIC_CONFIG_SCRIPT),
    )
    shutil.copy(
        os.path.join("windows", "scripts", NIC_CONFIG_SCRIPT),
        os.path.join(BUILD_DIR, NIC_CONFIG_SCRIPT),
    )
    shutil.copy(
        os.path.join("windows", "scripts", GEN_CERT_SCRIPT),
        os.path.join(BUILD_DIR, GEN_CERT_SCRIPT),
    )


def create_wxs_file():
    """
    Generate the .wxs file that is used by the WiX Toolset to generate the .msi file.

    This file contains the list of files to be included in the installer, as well as
    the UI configuration, and the custom actions to be executed.
    """
    wxs_content = f"""<?xml version="1.0" encoding="UTF-8"?>
<Wix xmlns="http://schemas.microsoft.com/wix/2006/wi" xmlns:util="http://schemas.microsoft.com/wix/UtilExtension">
    <Product
        Id="*"
        UpgradeCode="966ecd3d-a30f-4909-a0b4-0df045930e7d"
        Name="MassaStation"
        Language="1033"
        Version="{VERSION}"
        Manufacturer="MassaLabs"
    >
        <Package
            InstallerVersion="200"
            Compressed="yes"
            InstallScope="perMachine"
        />
        <MajorUpgrade
            Schedule="afterInstallInitialize"
            DowngradeErrorMessage="A later version of [ProductName] is already installed"
            AllowSameVersionUpgrades="yes"
        />

        <MediaTemplate EmbedCab="yes"/>
        
        <WixVariable Id="WixUILicenseRtf" Value="windows\\License.rtf" />
        <Property Id="MsiLogging" Value="voicewarmup!" />

        <Property Id="WIXUI_INSTALLDIR" Value="INSTALLDIR" />
        <Property Id="WIXUI_EXITDIALOGOPTIONALTEXT" Value="Launch MassaStation" />
        <Property Id="WIXUI_EXITDIALOGOPTIONALCHECKBOXTEXT" Value="Launch MassaStation" />
        <Property Id="WIXUI_EXITDIALOGOPTIONALCHECKBOX" Value="1" />

        <UI>
            <UIRef Id="WixUI_InstallDir" />
            <Publish Dialog="ExitDialog" Control="Finish" Event="DoAction" Value="LaunchMassaStation">WIXUI_EXITDIALOGOPTIONALCHECKBOX = 1</Publish>
        </UI>

        <Directory Id="TARGETDIR" Name="SourceDir">
            <Directory Id="ProgramFilesFolder">
                <Directory Id="INSTALLDIR" Name="MassaStation">
                    <Component Id="MassaStationServer" Guid="bc60f0be-065b-4738-968f-ce0e9b32bd01">
                        <File Id="MassaStationServerEXE" Name="{SERVER_BINARIES}" Source="{BUILD_DIR}\\{SERVER_BINARIES}" />
                        <File Id="MassaStationAppEXE" Name="{APP_BINARIES}" Source="{BUILD_DIR}\\{APP_BINARIES}" />
                        <File Id="AcrylicConfigScript" Name="{ACRYLIC_CONFIG_SCRIPT}" Source="{BUILD_DIR}\\{ACRYLIC_CONFIG_SCRIPT}" />
                        <File Id="NICConfigScript" Name="{NIC_CONFIG_SCRIPT}" Source="{BUILD_DIR}\\{NIC_CONFIG_SCRIPT}" />
                        <File Id="GenCertScript" Name="{GEN_CERT_SCRIPT}" Source="{BUILD_DIR}\\{GEN_CERT_SCRIPT}" />
                    </Component>
                    <Directory Id="MassaStationCerts" Name="certs">
                        <Component Id="CreateCertsDir" Guid="e96619b3-48a7-4629-8a19-e1c8270b331c">
                            <CreateFolder />
                        </Component>
                    </Directory>
                    <Directory Id="MassaStationPlugins" Name="plugins">
                        <Component Id="CreatePluginsDir" Guid="130fb4bb-cb51-4e28-a5e4-b7c58c846e02">
                            <CreateFolder>
                                <util:PermissionEx User="Users" GenericAll="yes"/>
                            </CreateFolder>
                        </Component>
                    </Directory>
                </Directory>
                <Directory Id="AcrylicDNSProxy" Name="Acrylic DNS Proxy">
                    <Component Id="Acrylic" Guid="563952aa-5f05-4c00-b3e0-6c004c36dc77">
                        <File Id="AcrylicZIP" Name="{ACRYLIC_ZIP}" Source="{BUILD_DIR}\\{ACRYLIC_ZIP}" />
                        <Condition><![CDATA[NOT Installed]]></Condition>
                    </Component>
                </Directory>
            </Directory>
            <Directory Id="TempDir" FileSource="[TempFolder]"></Directory>
        </Directory>

        <Feature Id="MainApplication" Title="Main Application" Level="1">
            <ComponentRef Id="MassaStationServer" />
            <ComponentRef Id="CreateCertsDir" />
            <ComponentRef Id="CreatePluginsDir" />
        </Feature>

        <Feature Id="Acrylic" Title="Acrylic DNS Proxy" Level="1">
            <ComponentRef Id="Acrylic" />
        </Feature>

        <CustomAction
            Id="LaunchMassaStation"
            Execute="immediate" 
            Return="asyncNoWait"
            Impersonate="yes"
            FileKey="MassaStationAppEXE"
            ExeCommand=""
        />

        <CustomAction Id="ExtractAcrylic"
            Directory="AcrylicDNSProxy"
            ExeCommand="powershell.exe -Command &quot;Expand-Archive '[AcrylicDNSProxy]{ACRYLIC_ZIP}' -d '[AcrylicDNSProxy]'&quot;"
            Execute="deferred"
            Impersonate="no"
            Return="check"
        />
        <CustomAction Id="DeleteAcrylicZip"
            Directory="AcrylicDNSProxy"
            ExeCommand="cmd /c del &quot;[AcrylicDNSProxy]{ACRYLIC_ZIP}&quot;"
            Execute="deferred"
            Impersonate="no"
            Return="check"
        />
        <CustomAction
            Id="InstallAcrylic"
            Directory="AcrylicDNSProxy"
            ExeCommand="cmd /c &quot;[AcrylicDNSProxy]InstallAcrylicService.bat&quot;"
            Execute="deferred"
            Impersonate="no"
            Return="check"
        />
        <CustomAction
            Id="ConfigureAcrylic"
            Directory="INSTALLDIR"
            ExeCommand="cmd /c &quot;[INSTALLDIR]{ACRYLIC_CONFIG_SCRIPT}&quot;"
            Execute="deferred"
            Impersonate="no"
            Return="check"
        />
        <CustomAction
            Id="ConfigureNetworkInterface"
            Directory="INSTALLDIR"
            ExeCommand="cmd /c &quot;[INSTALLDIR]{NIC_CONFIG_SCRIPT}&quot;"
            Execute="deferred"
            Impersonate="no"
            Return="check"
        />
        <CustomAction
            Id="GenerateCertificate"
            Directory="TempDir"
            ExeCommand="cmd /c &quot;[INSTALLDIR]{GEN_CERT_SCRIPT}&quot;"
            Execute="deferred"
            Impersonate="no"
            Return="check"
        />

        <InstallExecuteSequence>
            <Custom Action="ExtractAcrylic" Before="InstallAcrylic">NOT Installed AND NOT UpgradeCode</Custom>
            <Custom Action="InstallAcrylic" Before="ConfigureAcrylic">NOT Installed AND NOT UpgradeCode</Custom>
            <Custom Action="ConfigureAcrylic" Before="ConfigureNetworkInterface">NOT Installed</Custom>
            <Custom Action="ConfigureNetworkInterface" Before="GenerateCertificate">NOT Installed</Custom>
            <Custom Action="GenerateCertificate" Before="DeleteAcrylicZip">NOT Installed</Custom>
            <Custom Action="DeleteAcrylicZip" Before="InstallFinalize">NOT Installed</Custom>
        </InstallExecuteSequence>

    </Product>
</Wix>
"""
    with open(os.path.join(BUILD_DIR, "config.wxs"), "w", encoding="utf-8") as wxs_file:
        wxs_file.write(wxs_content)


def build_installer():
    """
    This function is the main function of this script.

    It downloads the binaries and builds the installer.
    """

    # URLs for downloading the binaries
    acrylic_url = "https://sourceforge.net/projects/acrylic/files/Acrylic/2.1.1/Acrylic-Portable.zip/download"
    massastation_server_url = "https://github.com/massalabs/thyra/releases/latest/download/thyra-server_windows_amd64"
    massastation_app_url = "https://github.com/massalabs/Thyra-Menu-Bar-App/releases/latest/download/ThyraApp_windows-amd64.exe"

    download_file(acrylic_url, ACRYLIC_ZIP)
    download_file(massastation_server_url, SERVER_BINARIES)
    download_file(massastation_app_url, APP_BINARIES)

    move_binaries()

    create_wxs_file()

    # Build installer
    try:
        subprocess.run(
            [
                os.path.join(WIX_DIR, "candle.exe"),
                "-ext",
                "WixUIExtension",
                "-ext",
                "WixUtilExtension",
                "-out",
                os.path.join(BUILD_DIR, "config.wixobj"),
                os.path.join(BUILD_DIR, "config.wxs"),
            ],
            check=True,
        )
        subprocess.run(
            [
                os.path.join(WIX_DIR, "light.exe"),
                "-ext",
                "WixUIExtension",
                "-ext",
                "WixUtilExtension",
                "-sval",
                "-out",
                f"massastation_{VERSION}_amd64.msi",
                os.path.join(BUILD_DIR, "config.wixobj"),
            ],
            check=True,
        )
    except subprocess.CalledProcessError as err:
        print("Error building installer: ", err)
        sys.exit(1)
    except err:
        print(err)
        sys.exit(1)


def install_dependencies():
    """
    Install dependencies if they are not already installed.

    The main dependency is the WiX Toolset.
    """

    if os.path.exists(WIX_DIR):
        return

    wixtoolset_url = (
        "https://wixdl.blob.core.windows.net/releases/v3.14.0.6526/wix314-binaries.zip"
    )

    download_file(wixtoolset_url, "wixtoolset.zip")

    os.mkdir(WIX_DIR)
    with zipfile.ZipFile("wixtoolset.zip", "r") as zip_ref:
        zip_ref.extractall(WIX_DIR)
    os.remove("wixtoolset.zip")


if __name__ == "__main__":
    # pylint: disable-next=pointless-string-statement
    """
    This is the main function of this script.

    It checks if the script is being run on Windows and if so,
    it installs the dependencies and builds the installer.
    """

    parser = argparse.ArgumentParser(
        description="Build the MassaStation installer for Windows."
    )
    parser.add_argument(
        "-v",
        "--version",
        action="store",
        dest="version",
        default=None,
        help="Set the version of the installer and MassaStation.",
    )
    parser.add_argument(
        "-f",
        "--force",
        action="store_true",
        dest="force_build",
        default=False,
        help="Force the build on non-Windows platforms.",
    )
    args = parser.parse_args()

    if args.version is None:
        print("Getting version from git tags")
        try:
            args.version = subprocess.run(
                ["git", "describe", "--tags", "--abbrev=0"],
                check=True,
                capture_output=True,
                text=True,
            ).stdout.strip()
            args.version = re.sub(r"^v", "", args.version)
        except subprocess.CalledProcessError as processErr:
            print("Error getting version: ", processErr)
            sys.exit(1)
        except Exception as gitErr:
            print("Error getting version: ", gitErr)
            sys.exit(1)

    VERSION = args.version

    if not sys.platform.startswith("win"):
        if not args.force_build:
            print(
                "This script is only meant to be run on Windows. If you wish to run it anyway, use the '-f' flag."
            )
            sys.exit(1)

    try:
        install_dependencies()
    except Exception as depErr:
        print("Error installing dependencies: ", depErr)
        sys.exit(1)
    try:
        build_installer()
    except Exception as buildErr:
        print("Error building installer: ", buildErr)
        sys.exit(1)
    sys.exit(0)
