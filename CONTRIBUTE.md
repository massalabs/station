# Contribute to Thyra
#### Clone the repo 
run 
```bash
git clone https://github.com/massalabs/thyra.git 
```

## Install dev dependencies 
 

#### Linux:

```bash
sudo apt update
sudo apt install -y build-essential libgl1-mesa-dev xorg-dev p7zip
```
#### all OS:

To develop on this project you will need :

- [go](https://go.dev/doc/install)
- [swagger](https://github.com/go-swagger/go-swagger) to generate go code from API documentation

#### build generated files:
#####

This command will build all generated files, including, go swagger code, Stringers.

Go swagger generate go code from API documentation, it can be installed using this command:

```bash
go install github.com/go-swagger/go-swagger/cmd/swagger@latest
```

launch code generation:

```bash
go generate
```



### From source 

#### Run the thyra-server
copy and paste the command line `thyra-server` in your terminal. You can now browse the websites on-chain seamlessly
```
go build -o thyra-server cmd/thyra-server/main.go
sudo setcap CAP_NET_BIND_SERVICE=+eip thyra-server 
./thyra-server 
```

### build front-end:

The frontend is built with react. Its source code needs to be compiled to be used locally.
You will need npm and node to be installed.

```bash
go generate ./web/...
```

- `go install github.com/go-swagger/go-swagger/cmd/swagger@latest`





