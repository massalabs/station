# Installation

* [Introduction](#introduction)
* [Step by step instructions](#step-by-step-instructions)
  * [MacOS](#macos)
    * [Automatically](#automatically)
    * [Manually](#manually)
      * [Thyra installation](#thyra-installation)
      * [DNS installation](#dns-installation)
  * [Linux](#linux)
    * [Automatically](#automatically)
    * [Manually](#manually)
      * [Thyra installation](#thyra-installation-1)
      * [DNS installation](#dns-installation-1)
  * [Windows](#windows)
    * [Automatically](#automatically)
    * [Manually](#manually)
      * [Thyra installation](#thyra-installation-2)
      * [DNS installation](#dns-installation-2)

## Introduction

This document will guide you through the installation process of the latest tagged version of Thyra.

> **_PREREQUISITES:_** Be comfortable with your terminal and have a recent version of MacOS, Windows or Linux.

> **_TROUBLESHOOTING:_** If you have trouble following this procedure, feel free to [open a question](https://github.com/massalabs/thyra/issues/new) describing your problem.

## Step by step instructions

Two steps are required to use our web on-chain product:

* Thyra installation : obtain the Thyra binary corresponding to your operating system (OS), rename it and make it executable.
* DNS configuration : install and configure your DNS to resolve the massa top level domain (*.massa) where Thyra runs.

You can either perform these two steps automatically, using an installation script, or manually using the binaries and command lines.

Now, let's move on to your OS section:

* [Linux](#linux)
* [Windows](#windows)
* [MacOS](#macos)

### MacOS

#### Automatically
> **_NOTE:_** If you're not a developer, it's better to use automatic installation

Simply use copy/paste the cmd line below in your terminal. The installation process will start and a success message will be displayed once done.


```sh
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/massalabs/thyra/main/scripts/macos_install.sh)"
```

Congratulation, your version of Thyra is now installed on your system and can be run by executing `thyra-server` in your terminal.
You can now browse the **websites on-chain** seamlessly. If you need to take the pressure off, maybe a little [game](http://flappy.massa) can help.
If you want to get down to business, you can start your [Massalian journey](http://my.massa/thyra/wallet) right away!




#### Manually

##### Thyra installation

Let's start by downloading the version of Thyra corresponding to your system:

* If you have an Intel CPU (amd64), you can download your executable [here](https://github.com/massalabs/thyra/releases/latest/download/thyra-server_darwin_amd64).
* If you have an Intel Apple Silicon CPU M1, M1 Pro, M2... (arm64), you can download your executable [here](https://github.com/massalabs/thyra/releases/latest/download/thyra-server_darwin_arm64).

Next, let's rename and make the downloaded binary executable by running the following command in your terminal:

```sh
mv thyra-server* thyra-server && chmod +x ./thyra-server
```

> **_NOTE:_** These commands should be executed directly from the directory where Thyra was downloaded.

Congratulation, your version of Thyra is now installed on your system and can be run by executing `./thyra-server` in your terminal.

> **_NOTE:_** If your DNS is already configured to handle the massa TLD, you're free to go. Otherwise, please follow the instructions in the next section.

##### DNS installation

> **_WARNING:_** If you already have a DNS service running that is not dnsmasq, you must configure it to redirect .massa to 127.0.0.1 (localhost).

> **_PREREQUISITE:_** Have `homebrew` already installed on your system. If not, you can follow the installation instructions [here](https://brew.sh).

Let's start by installing `dnsmasq`. This step can be safely skipped if it is already installed on your system.

```sh
brew install dnsmasq
```

Next, let's configure it to redirect .massa request to locahost:

```sh
echo 'address=/.massa/127.0.0.1' >> $(brew --prefix)/etc/dnsmasq.d/massa.conf
sudo mkdir -p /etc/resolver
sudo bash -c 'echo "nameserver 127.0.0.1" >> /etc/resolver/massa'

echo "Now we need you to type in your password to start the dnsmasq service."
sudo brew services start dnsmasq
```

Congratulations, you can now browse the **websites on-chain** seamlessly. If you need to take the pressure off, maybe a little [game](http://flappy.massa) can help.
If you want to get down to business, you can start your [Massalian journey](http://my.massa/thyra/wallet) right away!

### Linux

#### Automatically

Simply use copy/paste the cmd line below in your terminal. The installation process will start and a success message will be displayed once done.


```sh
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/massalabs/thyra/main/scripts/linux_install.sh)"
```

Congratulation, your version of Thyra is now installed on your system and can be run by executing `thyra-server` in your terminal.

You can now browse the **websites on-chain** seamlessly. If you need to take the pressure off, maybe a little [game](http://flappy.massa) can help.
If you want to get down to business, you can start your [Massalian journey](http://my.massa/thyra/wallet) right away!

> **_NOTE:_** Only Linux Ubuntu is currently supported.




#### Manually

##### Thyra installation

Let's start by downloading the version of Thyra corresponding to your system [here](https://github.com/massalabs/thyra/releases/latest/download/thyra-server_linux_amd64).

Next, let's rename and make the downloaded binary executable by running the following command in your terminal:

```sh
mv thyra-server* thyra-server && chmod +x ./thyra-server
```

> **_NOTE:_** These commands should be executed directly from the directory where Thyra was downloaded.

Congratulation, your version of Thyra is now installed on your system and can be run by executing `thyra-server` in your terminal.

> **_NOTE:_** If your DNS is already configured to handle the massa TLD, you're free to go. Otherwise, please follow the instructions in the next section.

##### DNS installation

> **_WARNING:_** If you already have a DNS service running that is not dnsmasq, you must configure it to redirect .massa to 127.0.0.1 (localhost).

###### dnsmasq (default)

> **_NOTE:_** If you have `NetworkManager` running, you must change its configuration to use `dnsmasq` as your local DNS. You can do this by running the following command:
>
>```sh
>sudo cp /etc/NetworkManager/NetworkManager.conf /etc/NetworkManager/NetworkManager.conf_backup_thyra_install && sudo sed -i "s/keyfile/keyfile\ndns=dnsmasq/g" /etc/NetworkManager/NetworkManager.conf
>```
>
> This command backs up `/etc/NetworkManager/NetworkManager.conf` file to `/etc/NetworkManager/NetworkManager.conf_backup_thyra_install`.

Then we must configure and restart the dnsmasq service:

```sh
sudo mkdir -p /etc/NetworkManager/dnsmasq.d/
echo "address=/.massa/127.0.0.1" | sudo tee -a /etc/NetworkManager/dnsmasq.d/massa.conf > /dev/null
sudo mv /etc/resolv.conf /etc/resolv.conf_backup_thyra_install && sudo ln -s /var/run/NetworkManager/resolv.conf /etc/resolv.conf
sudo systemctl restart NetworkManager
```

> **_NOTE:_** your `/etc/resolv.conf` file has been backed up to `/etc/resolv.conf_backup_thyra_install`.

Congratulations, you can now browse the **websites on-chain** seamlessly. If you need to take the pressure off, maybe a little [game](http://flappy.massa) can help.
If you want to get down to business, you can start your [Massalian journey](http://my.massa/thyra/wallet) right away!

### Windows

#### Automatically
> **_NOTE:_** If you're not a developer, it's better to use automatic installation

Follow the link below and download the file named "ThyraApp_windows-amd64.exe" on your computer. Then: 
1. Open it by double clicking
2. The script starts and here is what it does:
  * If you have Thyra already installed on your computer, it installs Thyra's icon tray.
  * If nothing is installed on your computer yet, it will install both.

![windows_icontray_V0](https://user-images.githubusercontent.com/109611779/212294116-05e1dd37-ed3f-4e3e-b034-b02d782bc4ee.png)

Congratulation, your version of Thyra is now installed on your system and you can "Start" your journey using the icon tray.

You can now browse the **websites on-chain** seamlessly. If you need to take the pressure off, maybe a little [game](http://flappy.massa) can help.
If you want to get down to business, you can start your [Massalian journey](http://my.massa/thyra/wallet) right away!


#### Manually

Let's start by downloading the version of Thyra corresponding to your system [here](https://github.com/massalabs/thyra/releases/latest/download/thyra-server_windows_amd64).

Next, you must rename it to `thyra-server.exe`.

Congratulation, your version of Thyra is now installed on your system and can be run by executing `thyra-server`.

> **_NOTE:_** If your DNS is already configured to handle the massa TLD, you're free to go. Otherwise, please follow the instructions in the next section.

##### DNS installation

> **_WARNING:_** If you already have a DNS service running that is not dnsmasq, you must configure it to redirect .massa to 127.0.0.1 (localhost).

> **_PREREQUISITE:_** Have `Acrylic` already installed on your system. If not, you can follow the installation instructions [here](https://mayakron.altervista.org/support/acrylic/Home.htm) and the OS configuration [here](https://mayakron.altervista.org/support/acrylic/Windows10Configuration.htm).

Let's start by configuring acrylic to redirect *.massa to locahost:

1. Open Acrylic config file: Open Acrylic DNS Proxy UI > File > Open Acrylic Hosts

2. Add `*.massa` top level domain to `AcrylicHosts.txt`:

```txt
127.0.0.1   *.massa
```

3. Save the file and reload Acrylic: Open Acrylic DNS Proxy UI > Actions > Restart Acrylic Service

Congratulations, you can now browse the **websites on-chain** seamlessly. If you need to take the pressure off, maybe a little [game](http://flappy.massa) can help.
If you want to get down to business, you can start your [Massalian journey](http://my.massa/thyra/wallet) right away!
