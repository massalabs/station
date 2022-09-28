# Installation

Currently, Thyra can be installed using prebuilt binaries. We might support other installation process later.

### MacOS

1. Download the MacOS binary corresponding to your system in the *asset* section of the [latest Release](https://github.com/massalabs/thyra/releases/latest).

Note: if you are using an Intel CPU, you need to download the `amd64` version, if you have an Apple Silicon CPU (M1, M1 Pro, M2...) you need to download the `arm64` version.

2. Install `dnsmasq` using `Homebrew`. If `Homebrew` isn't installed on your computer, install it using the instructions available [here](https://brew.sh).
* Open a terminal and run the following command:
```sh
brew install dnsmasq
```


From there, paste below cmd on your terminal.


3. Create `dnsmasq` config for `*.massa`

```sh
echo 'address=/.massa/127.0.0.1' >> $(brew --prefix)/etc/dnsmasq.d/massa.conf
```

4. Add your nameserver to resolvers

```sh
sudo mkdir -p /etc/resolver
sudo bash -c 'echo "nameserver 127.0.0.1" > /etc/resolver/massa'
```

5. Restart `dnsmasq`

```sh
sudo brew services start dnsmasq
```

6. Get in the directory you downloaded Thyra and run the following command:

For a Mac with an Intel CPU:
```sh
chmod +x ./thyra-server_darwin_amd64
```

And for a Mac with an Apple Silicon CPU:
```sh
chmod +x ./thyra-server_darwin_arm64
```

7. Start Thyra using the binary you downloaded in step 1

For a Mac with an Intel CPU, it should be:
```sh
./thyra-server_darwin_amd64
```

And for a Mac with an Apple Silicon CPU, it should be:
```sh
./thyra-server_darwin_arm64
```

You now should be able to access Thyra using your web browser at `http://my.massa/thyra/wallet/index.html`.
In case of error, feel free to [open an issue](https://github.com/massalabs/thyra/issues/new) describing your problem. 

### Linux

#### Disclaimer
For now, this installation process has been tested on Ubuntu 22.04 and might not work on other distribution. Feel free to [open an issue](https://github.com/massalabs/thyra/issues/new) describing your problem if you face any.

#### Ubuntu 

1. Download the Linux binary from the asset section of the [latest Release](https://github.com/massalabs/thyra/releases/latest).

2. Replace the used DNS by `NetworkManager` to use `dnsmasq`

```sh
sudo sed -i "s/keyfile/keyfile\ndns=dnsmasq/g" "/etc/NetworkManager/NetworkManager.conf"
```

3. Make sure `dnsmasq` configuration directory exist

```sh
sudo mkdir -p "/etc/NetworkManager/dnsmasq.d/"
```

4. Create `dnsmasq` config for `*.massa`

```sh
echo "address=/.massa/127.0.0.1" | sudo tee /etc/NetworkManager/dnsmasq.d/massa.conf > /dev/null
```

5. Update `/etc/resolv.conf` simlink to use `NetworkManager`

```sh
sudo rm "/etc/resolv.conf"
sudo ln -s "/var/run/NetworkManager/resolv.conf" "/etc/resolv.conf"
```

6. Restart NetworkManager

```sh
sudo systemctl restart NetworkManager
```

7. Get in the directory you downloaded Thyra and run the following command:

```sh
chmod +x ./thyra-server_linux_amd64
```

8. start Thyra using the binary you downloaded in step 1

```sh
./thyra-server_linux_amd64
```

If you get the following error: `listen tcp :80: bind: permission denied`, you might need to run the previous command as `sudo` or change the http port using the `-http-port` and https port using the `-https-port` arguments.

You now should be able to access Thyra using your web browser at `http://my.massa/thyra/wallet/index.html`.
In case of error, feel free to [open an issue](https://github.com/massalabs/thyra/issues/new) describing your problem. 

### Windows

1. Download the Windows binary from the asset section of the [latest Release](https://github.com/massalabs/thyra/releases/latest).

2. [Install Acrylic](https://mayakron.altervista.org/support/acrylic/Home.htm)

3. [Configure Acrylic DNS Proxy](https://mayakron.altervista.org/support/acrylic/Windows10Configuration.htm)

4. Open Acrylic config file: Open Acrylic DNS Proxy UI > File > Open Acrylic Hosts

5. Add `*.massa` top level domain to `AcrylicHosts.txt`: 
```txt
127.0.0.1   *.massa
```
Make sure to save using CTRL+S.

6. Restart Acrylic: Open Acrylic DNS Proxy UI > Actions > Restart Acrylic Service

7. Start Thyra by executing the binary you downloaded during step 1

You now should be able to access Thyra using your web browser at `http://my.massa/thyra/wallet/index.html`.
In case of error, feel free to [open an issue](https://github.com/massalabs/thyra/issues/new) describing your problem. 

