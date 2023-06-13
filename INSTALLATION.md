# Installation Guide

This guide provides step-by-step instructions for installing **MassaStation** on your computer. Please follow the instructions specific to your operating system.

* [Windows Installation](#windows-installation)
* [MacOS Installation](#macos-installation)
* [Debian Linux Installation](#debian-linux-installation)
* [Uninstallation](#uninstallation)
  * [MacOS](#macos)

> **Note:** **MassaStation** is currently available for Windows, MacOS, and Debian Linux. Support for other Linux distributions will be added in the future. Likewise, this application isn't working on virtual machines. It might be added in the future.

> **Troubleshooting:** If you encounter any issues during the installation process, do not hesitate to [open an issue](https://github.com/massalabs/thyra/issues/new) on GitHub.

## Windows Installation

1. Download the latest version of **MassaStation** installer for Windows (`.msi`) from the [official website](https://github.com/massalabs/thyra/releases/latest/download/massastation-installer_windows_amd64.msi).
2. Locate the downloaded `.msi` installer file and double-click on it to start the installation process.
3. Follow the on-screen instructions to proceed with the installation.
4. Once the installation is complete, you will see a confirmation message. Click "Finish" to exit the installer.

Congratulations! You have successfully installed **MassaStation** on your computer. You can start **MassaStation** by searching for it in the Start menu or by clicking on the icon on your desktop (if you chose to create one during the installation).

## MacOS Installation

1. Download the latest version of **MassaStation** installer for MacOS depending on your CPU architecture:

   * For Intel-based Macs (i5, i7, etc.), download the installer for `amd64` architecture from [here](https://github.com/massalabs/thyra/releases/latest/).
   * For Apple Silicon-based Macs (M1, M2, etc.), download the installer for `arm64` architecture from [here](https://github.com/massalabs/thyra/releases/latest/).

2. Locate the downloaded .pkg installer file and right-click on it.
3. From the context menu, select "Open" and then click "Open" again in the security pop-up window. This step is necessary because the installer is not signed by the App Store, and MacOS may block the installation by default.
4. If prompted, enter your administrator password to authorize the installation.
5. Follow the on-screen instructions to proceed with the installation.
6. Once the installation is complete, you will see a confirmation message. Click "Close" to exit the installer.

Congratulations! You have successfully installed **MassaStation** on your computer. You can start **MassaStation** by searching for it in the Applications folder or by clicking on the icon in the Launchpad.

## Debian Linux Installation

### GUI Installation

1. Download the latest version of **MassaStation** installer for Debian Linux (`.deb`)  from [here](https://github.com/massalabs/thyra/releases/latest/).
2. Open your file manager and navigate to the location where the `.deb` package is saved.
3. Right-click on the `.deb` package and choose "Open with Software Install" or a similar option.
4. The package manager will launch and display **MassaStation** installation page.
5. Review the package information and dependencies, if any, and click on the "Install" button.
6. If prompted, enter your administrator password to authorize the installation.
7. The installation will commence, and you will see a progress bar indicating the status.
8. Once the installation is complete, you will receive a notification confirming the successful installation.

### Terminal Installation using apt

1. Download the latest version of **MassaStation** installer for Debian Linux from [here](https://github.com/massalabs/thyra/releases/latest/).
2. Open a terminal on your Debian Linux system.
3. Navigate to the directory where the downloaded `.deb` package is located.
4. Run the following command to install the package:

   ```bash
   sudo apt install ./massastation-${{ VERSION }}_amd64.deb
   ```

5. Enter your administrator password when prompted and press Enter to confirm.
6. The installation will begin, and you will see progress information in the terminal.
7. Once the installation is complete, you can close the terminal.

Congratulations! You have successfully installed **MassaStation** on your computer. You can start **MassaStation** by searching for it in the Applications folder.

## Uninstallation

### Windows
To uninstall **MassaStation** from your Windows system, follow the steps below: 
1. Open your "Start" panel
2. Look for MassaStation application in your list
3. Click right on the application and click on uninstall

The application and all modules installed will be deleted from your desktop.

### MacOS

To uninstall **MassaStation** from your MacOS system, follow the steps below:

1. Open the Terminal application on your MacOS system.
2. Execute the following command in the terminal to download and run the **MassaStation** uninstaller script:

   ```bash
   /usr/local/share/massastation/uninstall.sh
   ```

   This command will remove MassaStation and its associated files from your system.
3. Follow any prompts or instructions provided by the uninstaller script. This may involve confirming the removal and providing your password for administrative privileges.
4. Once the uninstallation process is complete, you will receive a confirmation message indicating that **MassaStation** has been successfully uninstalled.

> **Note:** DNSMasq and Homebrew might have been installed on your system as dependencies for MassaStation. We do not remove these packages automatically as they may be used by other applications on your system.

### Linux

To uninstall **MassaStation** from your Linux system, follow the steps below:

1. Open the Terminal application on your Linux system.
2. Execute the following command in the terminal to download and run the **MassaStation** uninstaller script:

   ```bash
   sudo dpkg -r massastation
   ```

   This command will remove MassaStation and its associated files from your system.
3. Once the uninstallation process is complete, you will receive a confirmation message indicating that **MassaStation** has been successfully uninstalled.
