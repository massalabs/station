# Installation

* [Introduction](#introduction)
* [Step by step instructions](#step-by-step-instructions)
  * [MacOS](#macos)
  * [Linux](#linux)
  * [Windows](#windows)

## Introduction

This document will guide you through the installation process of the latest tagged version of Thyra.

> **_PREREQUISITES:_** Be comfortable with your terminal and have a recent version of MacOS, Windows or Linux.

> **_TROUBLESHOOTING:_** If you have trouble following this procedure, feel free to [open a question](https://github.com/massalabs/thyra/issues/new) describing your problem.

## Step by step instructions

Two steps are required to use our web on-chain product:

* Thyra installation : obtain the Thyra binary corresponding to your operating system (OS), rename it and make it executable.
* DNS configuration : install and configure your DNS to resolve the massa top level domain (*.massa) where Thyra runs.

The following installation guide allows you to perform these steps automatically using a simple script to run on your terminal.

Now, let's move on to your OS section:

* [Linux](#linux)
* [Windows](#windows)
* [MacOS](#macos)


### MacOS

Simply use copy/paste the cmd line below in your terminal. The installation process will start and a success message will be displayed once done.


```sh
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/massalabs/thyra/main/scripts/macos_install.sh)"
```

Congratulation, your version of Thyra is now installed on your system and can be run by executing `thyra-server` in your terminal.
You can now browse the **websites on-chain** seamlessly. If you need to take the pressure off, maybe a little [game](http://flappy.massa) can help.
If you want to get down to business, you can start your [Massalian journey](http://my.massa/thyra/wallet) right away!



### Linux


Simply use copy/paste the cmd line below in your terminal. The installation process will start and a success message will be displayed once done.


```sh
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/massalabs/thyra/main/scripts/linux_install.sh)"
```

Congratulation, your version of Thyra is now installed on your system and can be run by executing `thyra-server` in your terminal.

You can now browse the **websites on-chain** seamlessly. If you need to take the pressure off, maybe a little [game](http://flappy.massa) can help.
If you want to get down to business, you can start your [Massalian journey](http://my.massa/thyra/wallet) right away!

> **_NOTE:_** Only Linux Ubuntu is currently supported.


### Windows


Follow [this link](https://github.com/massalabs/Thyra-Menu-Bar-App/releases/tag/v0.0.1) and download the file named "ThyraApp_windows-amd64.exe" on your computer. Then: 
1. Open it by double clicking
2. The script starts and here is what it does:
  * If you have Thyra already installed on your computer, it installs Thyra's icon tray.
  * If nothing is installed on your computer yet, it will install both.

![windows_icontray_V0](https://user-images.githubusercontent.com/109611779/212294116-05e1dd37-ed3f-4e3e-b034-b02d782bc4ee.png)

Congratulation, your version of Thyra is now installed on your system and you can "Start" your journey using the icon tray.

You can now browse the **websites on-chain** seamlessly. If you need to take the pressure off, maybe a little [game](http://flappy.massa) can help.
If you want to get down to business, you can start your [Massalian journey](http://my.massa/thyra/wallet) right away!


