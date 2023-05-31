# Installation Guide

This guide provides step-by-step instructions for installing **MassaStation** on your computer. Please follow the instructions specific to your operating system.

* [Windows Installation](#windows-installation)
* [MacOS Installation](#macos-installation)
* [Debian Linux Installation](#debian-linux-installation)

> **Note:** **MassaStation** is currently available for Windows, MacOS, and Debian Linux. Support for other Linux distributions will be added in the future.

> **Troubleshooting:** If you encounter any issues during the installation process, do not hesitate to [open an issue](https://github.com/massalabs/thyra/issues/new) on GitHub.

## Windows Installation

1. Download the latest version of **MassaStation** installer for Windows from the [official website](https://github.com/massalabs/thyra/releases/latest/download/massastation-installer_windows_amd64.msi).
2. Locate the downloaded `.msi` installer file and double-click on it to start the installation process.
3. Follow the on-screen instructions to proceed with the installation.
4. Once the installation is complete, you will see a confirmation message. Click "Finish" to exit the installer.

Congratulations! You have successfully installed **MassaStation** on your computer. You can start MassaStation by searching for it in the Start menu or by clicking on the icon on your desktop (if you chose to create one during the installation).

## MacOS Installation

1. Download the latest version of **MassaStation** installer for MacOS depending on your CPU architecture:

   - For Intel-based Macs (i5, i7, etc.), download the installer for `amd64` architecture from [here](https://github.com/massalabs/thyra/releases/latest/).
   - For Apple Silicon-based Macs (M1, M2, etc.), download the installer for `arm64` architecture from [here](https://github.com/massalabs/thyra/releases/latest/).

2. Locate the downloaded .pkg installer file and right-click on it.
3. From the context menu, select "Open" and then click "Open" again in the security pop-up window. This step is necessary because the installer is not signed by the App Store, and MacOS may block the installation by default.
4. If prompted, enter your administrator password to authorize the installation.
5. Follow the on-screen instructions to proceed with the installation.
6. Once the installation is complete, you will see a confirmation message. Click "Close" to exit the installer.

Congratulations! You have successfully installed **MassaStation** on your computer. You can start MassaStation by searching for it in the Applications folder or by clicking on the icon in the Launchpad.

## Debian Linux Installation

### GUI Installation

1. Download the latest version of **MassaStation** installer for Debian Linux from [here](https://github.com/massalabs/thyra/releases/latest/).
2. Open your file manager and navigate to the location where the `.deb` package is saved.
3. Right-click on the `.deb` package and choose "Open with Software Install" or a similar option.
4. The package manager will launch and display **MassaStation** installation page.
5. Review the package information and dependencies, if any, and click on the "Install" button.
6. If prompted, enter your administrator password to authorize the installation.
7. The installation will commence, and you will see a progress bar indicating the status.
8. Once the installation is complete, you will receive a notification confirming the successful installation.

### Terminal Installation using apt

1. Open a terminal on your Debian Linux system.
2. Navigate to the directory where the downloaded `.deb` package is located.
3. Run the following command to install the package:

   ```
   sudo apt install ./massastation-${{ VERSION }}_amd64.deb
   ```

4. Enter your administrator password when prompted and press Enter to confirm.
5. The installation will begin, and you will see progress information in the terminal.
6. Once the installation is complete, you can close the terminal.

Congratulations! You have successfully installed **MassaStation** on your computer. You can start MassaStation by searching for it in the Applications folder.
