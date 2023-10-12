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

# Metadata
MANUFACTURER = "MassaLabs"
PRODUCT_NAME = "MassaStation"
VERSION = "0.0.0"

# Binaries to be included in the installer
MASSASTATION_BINARY = "massastation.exe"
ACRYLIC_ZIP = "acrylic.zip"
MARTOOLS_ZIP = "martools.zip"
WIXTOOLSET_ZIP = "wixtoolset.zip"

# Scripts to be included in the installer
ADD_STATION_TO_HOSTS_SCRIPT = "add_station_to_hosts.bat"
ACRYLIC_CONFIG_SCRIPT = "configure_acrylic.bat"
NIC_CONFIG_SCRIPT = "configure_network_interfaces.bat"
NIC_RESET_SCRIPT = "reset_network_interfaces.bat"
RUN_VBS = "run.vbs"
LOGO = "logo.ico"
WIX_DIR = "wixtoolset"

# URLs to download Acrylic DNS Proxy and the WiX Toolset
ACRYLIC_URL = "https://sourceforge.net/projects/acrylic/files/Acrylic/2.1.1/Acrylic-Portable.zip/download"
MARTOOLS_URL = "https://archive.torproject.org/tor-package-archive/torbrowser/12.5.1/mar-tools-win64.zip"
WIXTOOLSET_URL = (
    "https://wixdl.blob.core.windows.net/releases/v3.14.0.6526/wix314-binaries.zip"
)

def download_file(url, filename):
    """
    Download a given file from a given URL.
    """
    urllib.request.urlretrieve(url, filename)


def install_massastation_build_dependencies():
    """
    Install the dependencies required to build massastation.
    """
    # Install Go dependencies
    subprocess.run(
        ["go", "install", "github.com/go-swagger/go-swagger/cmd/swagger@latest"],
        check=True,
    )
    subprocess.run(
        ["go", "install", "golang.org/x/tools/cmd/stringer@latest"], check=True
    )
    subprocess.run(["go", "install", "fyne.io/fyne/v2/cmd/fyne@latest"], check=True)


def build_massastation():
    """
    Build the MassaStation binary from source.
    """
    install_massastation_build_dependencies()

    subprocess.run(["go", "generate", "../..."], check=True)
    os.environ["CGO_ENABLED"] = "1"

    subprocess.run([
        "go",
        "build", 
        "-ldflags",
        f"-X github.com/massalabs/station/int/config.Version={VERSION}",
        "-o",
        "massastation.exe",
        "../cmd/massastation/"
    ], check=True)


def move_binaries():
    """
    Move the binaries and scripts to the build directory.

    If the build directory already exists, it is deleted first.
    """
    if os.path.exists(BUILD_DIR):
        shutil.rmtree(BUILD_DIR)

    os.makedirs(BUILD_DIR)

    shutil.copy(
        os.path.join(MASSASTATION_BINARY),
        os.path.join(BUILD_DIR, MASSASTATION_BINARY),
    )
    shutil.copy(ACRYLIC_ZIP, os.path.join(BUILD_DIR, ACRYLIC_ZIP))
    shutil.copy(MARTOOLS_ZIP, os.path.join(BUILD_DIR, MARTOOLS_ZIP))

    shutil.copy(
        os.path.join("windows", "scripts", ADD_STATION_TO_HOSTS_SCRIPT),
        os.path.join(BUILD_DIR, ADD_STATION_TO_HOSTS_SCRIPT),
    )
    shutil.copy(
        os.path.join("windows", "scripts", ACRYLIC_CONFIG_SCRIPT),
        os.path.join(BUILD_DIR, ACRYLIC_CONFIG_SCRIPT),
    )
    shutil.copy(
        os.path.join("windows", "scripts", NIC_CONFIG_SCRIPT),
        os.path.join(BUILD_DIR, NIC_CONFIG_SCRIPT),
    )
    shutil.copy(
        os.path.join("windows", "scripts", NIC_RESET_SCRIPT),
        os.path.join(BUILD_DIR, NIC_RESET_SCRIPT),
    )
    shutil.copy(
        os.path.join("windows", "scripts", RUN_VBS),
        os.path.join(BUILD_DIR, RUN_VBS),
    )
    shutil.copy(
        os.path.join("windows","assets", LOGO),
        os.path.join(BUILD_DIR, LOGO),
    )


def create_wxs_file():
    """
    Generate the .wxs file that is used by the WiX Toolset to generate the .msi file.

    This file contains the list of files to be included in the installer, as well as
    the UI configuration, and the custom actions to be executed.
    """
    installer_logfile = f"$env:Temp\\MASSASTATION_INSTALLER_{VERSION}.log"

    wxs_content = f"""<?xml version="1.0" encoding="UTF-8"?>
<Wix xmlns="http://schemas.microsoft.com/wix/2006/wi" xmlns:util="http://schemas.microsoft.com/wix/UtilExtension">
    <Product
        Id="*"
        UpgradeCode="966ecd3d-a30f-4909-a0b4-0df045930e7d"
        Name="{PRODUCT_NAME}"
        Language="1033"
        Version="{VERSION}"
        Manufacturer="{MANUFACTURER}"
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

        <Property Id="WINDOWSEDITION" Secure="yes">
            <RegistrySearch Id="WindowsEditionReg" Root="HKLM" Key="SOFTWARE\Microsoft\Windows NT\CurrentVersion" Name="CurrentBuildNumber" Type="raw" />
        </Property>

        <Condition Message="Oops! Massa Station cannot be installed on your current Windows version ([WINDOWSEDITION]).\nWe require Windows 10 or a newer version for compatibility. Please consider updating your operating system to enjoy Massa Station.">
            <!--
                We require Windows 10 or a newer version for compatibility. First version of Windows 10 is 10240.
                We don't want to allow installation on Server versions of Windows.
            -->
            <![CDATA[Installed OR (WINDOWSEDITION >= 10240 AND MsiNTProductType = 1)]]>
        </Condition>

        <MediaTemplate EmbedCab="yes"/>
        <!--
            We don't need a license agreement for now.
            <WixVariable Id="WixUILicenseRtf" Value="windows\\License.rtf" />
        -->
        <Property Id="MsiLogging" Value="voicewarmup!" />

        <Property Id="WIXUI_INSTALLDIR" Value="INSTALLDIR" />
        <Property Id="WIXUI_EXITDIALOGOPTIONALTEXT" Value="Launch MassaStation" />
        <Property Id="WIXUI_EXITDIALOGOPTIONALCHECKBOXTEXT" Value="Launch MassaStation" />
        <Property Id="WIXUI_EXITDIALOGOPTIONALCHECKBOX" Value="1" />

        <UI>
            <UIRef Id="WixUI_Mondo" />
            <Publish Dialog="WelcomeDlg" Control="Next" Event="NewDialog" Value="SetupTypeDlg" Order="3">1</Publish>
            <Publish Dialog="SetupTypeDlg" Control="Back" Event="NewDialog" Value="WelcomeDlg" Order="3">1</Publish>
            <Publish Dialog="ExitDialog" Control="Finish" Event="DoAction" Value="LaunchMassaStation">WIXUI_EXITDIALOGOPTIONALCHECKBOX = 1</Publish>
        </UI>

        <Directory Id="TARGETDIR" Name="SourceDir">
            <Directory Id="ProgramFilesFolder">
                <Directory Id="INSTALLDIR" Name="MassaStation">
                    <Component Id="MassaStationServer" Guid="bc60f0be-065b-4738-968f-ce0e9b32bd01">
                        <File Id="MassaStationAppEXE" Name="{MASSASTATION_BINARY}" Source="{BUILD_DIR}\\{MASSASTATION_BINARY}" />
                        <File Id="AddStationToHostsScript" Name="{ADD_STATION_TO_HOSTS_SCRIPT}" Source="{BUILD_DIR}\\{ADD_STATION_TO_HOSTS_SCRIPT}" />
                        <File Id="AcrylicConfigScript" Name="{ACRYLIC_CONFIG_SCRIPT}" Source="{BUILD_DIR}\\{ACRYLIC_CONFIG_SCRIPT}" />
                        <File Id="NICConfigScript" Name="{NIC_CONFIG_SCRIPT}" Source="{BUILD_DIR}\\{NIC_CONFIG_SCRIPT}" />
                        <File Id="NICResetScript" Name="{NIC_RESET_SCRIPT}" Source="{BUILD_DIR}\\{NIC_RESET_SCRIPT}" />
                        <File Id="MassaStationRunScript" Name="{RUN_VBS}" Source="{BUILD_DIR}\\{RUN_VBS}" />
                    </Component>
                    <Directory Id="MassaStationCerts" Name="certs">
                        <Component Id="CreateCertsDir" Guid="e96619b3-48a7-4629-8a19-e1c8270b331c">
                            <CreateFolder>
                                <util:PermissionEx User="Users" GenericAll="yes"/>
                            </CreateFolder>
                        </Component>
                    </Directory>
                    <Directory Id="MassaStationLogs" Name="logs">
                        <Component Id="CreateLogsDir" Guid="4b24bfe1-c564-47a7-95d5-1268c661ef8a">
                            <CreateFolder>
                                <util:PermissionEx User="Users" GenericAll="yes"/>
                            </CreateFolder>
                        </Component>
                    </Directory>
                    <Directory Id="MassaStationMartools" Name="mar-tools">
                        <Component Id="CreateMartoolsDir" Guid="ac6a3582-6996-4d79-80c2-1cd2dc462da8">
                            <File Id="MartoolsZIP" Name="{MARTOOLS_ZIP}" Source="{BUILD_DIR}\\{MARTOOLS_ZIP}" />
                        </Component>
                    </Directory>
                    <Directory Id="MassaStationPlugins" Name="plugins">
                        <Component Id="CreatePluginsDir" Guid="130fb4bb-cb51-4e28-a5e4-b7c58c846e02">
                            <CreateFolder>
                                <util:PermissionEx User="Users" GenericAll="yes"/>
                            </CreateFolder>
                        </Component>
                    </Directory>
                    <Directory Id="MassaStationWebsitesCache" Name="websitesCache">
                        <Component Id="CreateWebsitesCacheDir" Guid="3932ae3c-1755-4612-8c02-5a8e1ee2c531">
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

            <Directory Id="ProgramMenuFolder" Name="Programs">
                <Directory Id="ApplicationProgramsFolder" Name="{MANUFACTURER}">
                    <Component Id="ApplicationShortcutProgramMenu" Guid="e2f5b2a0-0b0a-4b1e-9b0e-9b0e9b0e9b0e">
                        <Shortcut
                            Id="ApplicationStartMenuShortcut"
                            Name="{PRODUCT_NAME}"
                            Target="[#MassaStationRunScript]"
                            WorkingDirectory="INSTALLDIR"
                            Icon ="MassaStationIconProgramMenu"
                        >
                            <Icon Id="MassaStationIconProgramMenu" SourceFile="{BUILD_DIR}\\{LOGO}" />
                        </Shortcut>
                        <RemoveFolder Id="ApplicationProgramsFolder" On="uninstall" />
                        <RegistryValue Root="HKCU" Key="Software\{MANUFACTURER}\{PRODUCT_NAME}" Name="installed" Type="integer" Value="1" KeyPath="yes" />
                    </Component>
                </Directory>
            </Directory>

            <Directory Id="DesktopFolder" Name="Desktop">
                <Component Id="ApplicationShortcutDesktop" Guid="3e6f0b0e-1e0b-5a3c-7b0c-9c007a32f0e9">
                    <Shortcut Id="ApplicationDesktopShortcut"
                        Name="{PRODUCT_NAME}"
                        Target="[#MassaStationRunScript]"
                        WorkingDirectory="INSTALLDIR"
                        Icon="MassaStationIconDesktop"
                    >
                        <Icon Id="MassaStationIconDesktop" SourceFile="{BUILD_DIR}\\{LOGO}" />
                    </Shortcut>
                </Component>
            </Directory>
        </Directory>

        <Feature Id="MassaStation" Title="Massa Station" Description="Your gateway to the decentralized web" Level="1" Absent="disallow" AllowAdvertise="no" ConfigurableDirectory="INSTALLDIR">
            <ComponentRef Id="MassaStationServer" />
            <ComponentRef Id="CreateCertsDir" />
            <ComponentRef Id="CreatePluginsDir" />
            <ComponentRef Id="CreateLogsDir" />
            <ComponentRef Id="CreateWebsitesCacheDir" />
            <ComponentRef Id="CreateMartoolsDir" />
        </Feature>

        <Feature Id="DesktopShortcut" Title="Desktop Shortcut" Level="1" Absent="disallow" AllowAdvertise="no">
            <ComponentRef Id="ApplicationShortcutDesktop" />
        </Feature>

        <Feature Id="ProgramMenuShortcut" Title="Program Menu Shortcut" Level="1" Absent="disallow" AllowAdvertise="no">
            <ComponentRef Id="ApplicationShortcutProgramMenu" />
        </Feature>

        <Feature Id="AcrylicDNS" Title="DNS" Description="A DNS server that can be used to resolve .massa domains." Level="10">
            <ComponentRef Id="Acrylic" />
        </Feature>

        <CustomAction
            Id='LaunchMassaStation'
            Directory="INSTALLDIR"
            ExeCommand="cmd /c &quot;[INSTALLDIR]{RUN_VBS}&quot;"
            Execute='immediate'
            Return='ignore'
        />

        <SetProperty Id="ExtractMartools" Value="&quot;powershell.exe&quot; -Command &quot;Expand-Archive '[MassaStationMartools]{MARTOOLS_ZIP}' -d '[INSTALLDIR]'&quot;" Before="ExtractMartools" Sequence="execute" />
        <CustomAction
            Id="ExtractMartools"
            BinaryKey="WixCA"
            DllEntry="CAQuietExec"
            Execute="deferred"
            Impersonate="no"
            Return="check"
        />

        <SetProperty Id="DeleteMartoolsZip" Value="&quot;cmd.exe&quot; /c del &quot;[MassaStationMartools]{MARTOOLS_ZIP}&quot;" Before="DeleteMartoolsZip" Sequence="execute" />
        <CustomAction
            Id="DeleteMartoolsZip"
            BinaryKey="WixCA"
            DllEntry="CAQuietExec"
            Execute="deferred"
            Impersonate="no"
            Return="check"
        />

        <SetProperty Id="ExtractAcrylic" Value="&quot;powershell.exe&quot; -Command &quot;Expand-Archive '[AcrylicDNSProxy]{ACRYLIC_ZIP}' -d '[AcrylicDNSProxy]'&quot;" Before="ExtractAcrylic" Sequence="execute" />
        <CustomAction
            Id="ExtractAcrylic"
            BinaryKey="WixCA"
            DllEntry="CAQuietExec"
            Execute="deferred"
            Impersonate="no"
            Return="check"
        />

        <SetProperty Id="DeleteAcrylicZip" Value="&quot;cmd.exe&quot; /c del &quot;[AcrylicDNSProxy]{ACRYLIC_ZIP}&quot;" Before="DeleteAcrylicZip" Sequence="execute" />
        <CustomAction
            Id="DeleteAcrylicZip"
            BinaryKey="WixCA"
            DllEntry="CAQuietExec"
            Execute="deferred"
            Impersonate="no"
            Return="check"
        />

        <SetProperty Id="InstallAcrylic" Value="&quot;powershell.exe&quot; -Command &quot; Set-Location '[AcrylicDNSProxy]';  &amp; '[AcrylicDNSProxy]InstallAcrylicService.bat'&quot; >> {installer_logfile} 2>&amp;1" Before="InstallAcrylic" Sequence="execute" />
        <CustomAction
            Id="InstallAcrylic"
            BinaryKey="WixCA"
            DllEntry="CAQuietExec"
            Execute="deferred"
            Impersonate="no"
            Return="check"
        />

        <SetProperty Id="RollbackAcrylicInstall" Value="&quot;powershell.exe&quot; -Command &quot; Set-Location '[AcrylicDNSProxy]'; &amp; '[AcrylicDNSProxy]UninstallAcrylicService.bat'&quot; >> {installer_logfile} 2>&amp;1" Before="RollbackAcrylicInstall" Sequence="execute" />
        <CustomAction
            Id="RollbackAcrylicInstall"
            BinaryKey="WixCA"
            DllEntry="CAQuietExec"
            Execute="rollback"
            Impersonate="no"
            Return="ignore"
        />

        <SetProperty Id="AddStationToHosts" Value="&quot;powershell.exe&quot; -Command &quot; &amp; '[INSTALLDIR]{ADD_STATION_TO_HOSTS_SCRIPT}'&quot; >> {installer_logfile} 2>&amp;1" Before="AddStationToHosts" Sequence="execute" />
        <CustomAction
            Id="AddStationToHosts"
            BinaryKey="WixCA"
            DllEntry="CAQuietExec"
            Execute="deferred"
            Impersonate="no"
            Return="check"
        />

        <SetProperty Id="ConfigureAcrylic" Value="&quot;powershell.exe&quot; -Command &quot; &amp; '[INSTALLDIR]{ACRYLIC_CONFIG_SCRIPT}' '[AcrylicDNSProxy]' &quot; >> {installer_logfile} 2>&amp;1" Before="ConfigureAcrylic" Sequence="execute" />
        <CustomAction
            Id="ConfigureAcrylic"
            BinaryKey="WixCA"
            DllEntry="CAQuietExec"
            Execute="deferred"
            Impersonate="no"
            Return="check"
        />

        <SetProperty Id="ConfigureNetworkInterface" Value="&quot;powershell.exe&quot; -Command &quot; &amp; '[INSTALLDIR]{NIC_CONFIG_SCRIPT}'&quot; >> {installer_logfile} 2>&amp;1" Before="ConfigureNetworkInterface" Sequence="execute" />
        <CustomAction
            Id="ConfigureNetworkInterface"
            BinaryKey="WixCA"
            DllEntry="CAQuietExec"
            Execute="deferred"
            Impersonate="no"
            Return="check"
        />

        <SetProperty Id="RollbackNetworkInterface" Value="&quot;powershell.exe&quot; -Command &quot; &amp; '[INSTALLDIR]{NIC_RESET_SCRIPT}'&quot; >> {installer_logfile} 2>&amp;1" Before="ResetNetworkInterface" Sequence="execute" />
        <CustomAction
            Id="RollbackNetworkInterface"
            BinaryKey="WixCA"
            DllEntry="CAQuietExec"
            Execute="rollback"
            Impersonate="no"
            Return="ignore"
        />

        <SetProperty Id="UninstallAcrylic" Value="&quot;powershell.exe&quot; -Command &quot; Set-Location '[AcrylicDNSProxy]'; &amp; '[AcrylicDNSProxy]UninstallAcrylicService.bat'&quot; >> {installer_logfile} 2>&amp;1" Before="UninstallAcrylic" Sequence="execute" />
        <CustomAction
            Id="UninstallAcrylic"
            BinaryKey="WixCA"
            DllEntry="CAQuietExec"
            Execute="deferred"
            Impersonate="no"
            Return="ignore"
        />

        <SetProperty Id="ResetNetworkInterface" Value="&quot;powershell.exe&quot; -Command &quot; &amp; '[INSTALLDIR]{NIC_RESET_SCRIPT}'&quot; >> {installer_logfile} 2>&amp;1" Before="ResetNetworkInterface" Sequence="execute" />
        <CustomAction
            Id="ResetNetworkInterface"
            BinaryKey="WixCA"
            DllEntry="CAQuietExec"
            Execute="deferred"
            Impersonate="no"
            Return="check"
        />

        <SetProperty Id="RemoveInstallDir" Value="&quot;cmd.exe&quot; /c rmdir /s /q &quot;[INSTALLDIR]&quot;" Before="RemoveInstallDir" Sequence="execute" />
        <CustomAction
            Id="RemoveInstallDir"
            BinaryKey="WixCA"
            DllEntry="CAQuietExec"
            Execute="deferred"
            Impersonate="no"
            Return="ignore"
        />

        <SetProperty Id="RemoveAcrylicDir" Value="&quot;cmd.exe&quot; /c rmdir /s /q &quot;[AcrylicDNSProxy]&quot;" Before="RemoveAcrylicDir" Sequence="execute" />
        <CustomAction
            Id="RemoveAcrylicDir"
            BinaryKey="WixCA"
            DllEntry="CAQuietExec"
            Execute="deferred"
            Impersonate="no"
            Return="ignore"
        />

        <!-- This Custom Action removes the line containing "station.massa" from the hosts file. -->
        <!-- We can't directly overwrite the hosts file because it is protected by Windows. But creating a new file and moving it to the hosts file works. -->
        <SetProperty Id="RemoveStationFromHosts" Value="&quot;cmd.exe&quot; /c findstr /v station.massa %windir%\System32\drivers\etc\hosts > %windir%\System32\drivers\etc\hosts.new &amp; move /y %windir%\System32\drivers\etc\hosts.new %windir%\System32\drivers\etc\hosts" Before="RemoveStationFromHosts" Sequence="execute" />
        <CustomAction
            Id="RemoveStationFromHosts"
            BinaryKey="WixCA"
            DllEntry="CAQuietExec"
            Execute="deferred"
            Impersonate="no"
            Return="ignore"
        />

        <InstallExecuteSequence>
            <Custom Action="ExtractMartools" Before="DeleteMartoolsZip">NOT Installed</Custom>
            <Custom Action="ExtractAcrylic" Before="InstallAcrylic">NOT Installed AND <![CDATA[&AcrylicDNS=3]]></Custom>
            <Custom Action="InstallAcrylic" Before="ConfigureAcrylic">NOT Installed AND <![CDATA[&AcrylicDNS=3]]></Custom>
            <Custom Action="RollbackAcrylicInstall" Before="InstallAcrylic">NOT Installed AND <![CDATA[&AcrylicDNS=3]]></Custom>
            <Custom Action="ConfigureAcrylic" Before="ConfigureNetworkInterface">NOT Installed AND <![CDATA[&AcrylicDNS=3]]></Custom>
            <Custom Action="ConfigureNetworkInterface" Before="DeleteAcrylicZip">NOT Installed AND <![CDATA[&AcrylicDNS=3]]></Custom>
            <Custom Action="RollbackNetworkInterface" Before="ConfigureNetworkInterface">NOT Installed AND <![CDATA[&AcrylicDNS=3]]></Custom>
            <Custom Action="AddStationToHosts" Before="InstallFinalize">NOT Installed</Custom>
            <Custom Action="DeleteAcrylicZip" Before="InstallFinalize">NOT Installed AND <![CDATA[&AcrylicDNS=3]]></Custom>
            <Custom Action="DeleteMartoolsZip" Before="InstallFinalize">NOT Installed</Custom>

            <Custom Action='UninstallAcrylic' Before='RemoveFiles'>REMOVE="ALL" AND NOT UPGRADINGPRODUCTCODE AND <![CDATA[&AcrylicDNS=2]]></Custom>
            <Custom Action='ResetNetworkInterface' Before='RemoveFiles'>REMOVE="ALL" AND NOT UPGRADINGPRODUCTCODE AND <![CDATA[&AcrylicDNS=2]]></Custom>
            <Custom Action='RemoveStationFromHosts' Before='RemoveFiles'>REMOVE="ALL" AND NOT UPGRADINGPRODUCTCODE</Custom>
            <Custom Action='RemoveInstallDir' After='RemoveFiles'>REMOVE="ALL" AND NOT UPGRADINGPRODUCTCODE</Custom>
            <Custom Action='RemoveAcrylicDir' After='RemoveFiles'>REMOVE="ALL" AND NOT UPGRADINGPRODUCTCODE AND <![CDATA[&AcrylicDNS=2]]></Custom>
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

    download_file(ACRYLIC_URL, ACRYLIC_ZIP)
    download_file(MARTOOLS_URL, MARTOOLS_ZIP)

    if not os.path.exists(MASSASTATION_BINARY):
        build_massastation()

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
    print(f"Installer built successfully: massastation_{VERSION}_amd64.msi")


def install_dependencies():
    """
    Install dependencies if they are not already installed.

    This function installs WixToolset.
    """

    if not os.path.exists(WIX_DIR):
        download_file(WIXTOOLSET_URL, WIXTOOLSET_ZIP)

        os.mkdir(WIX_DIR)
        with zipfile.ZipFile(WIXTOOLSET_ZIP, "r") as zip_ref:
            zip_ref.extractall(WIX_DIR)
        os.remove(WIXTOOLSET_ZIP)


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
        "-f",
        "--force",
        action="store_true",
        dest="force_build",
        default=False,
        help="Force the build on non-Windows platforms.",
    )
    args = parser.parse_args()

    # Get $VERSION from environment variable
    VERSION = os.environ.get("VERSION")
    if VERSION is None or VERSION == "":
        print("Getting version from git tags")
        try:
            VERSION = subprocess.run(
                ["git", "describe", "--tags", "--abbrev=0"],
                check=True,
                capture_output=True,
                text=True,
            ).stdout.strip()
            VERSION = re.sub(r"^v", "", VERSION)
        except subprocess.CalledProcessError as processErr:
            print("Error getting version: ", processErr)
            sys.exit(1)
        except Exception as gitErr:
            print("Error getting version: ", gitErr)
            sys.exit(1)
    else:
        # Remove the "v" from the version if it exists
        VERSION = re.sub(r"^v", "", VERSION)

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
