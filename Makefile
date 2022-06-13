# In this project `make` is only used to generate native libraries to be bridged in a non go language.
# If you simply want to use an application, you're looking at the wrong file :-).

.PHONY: all, linux_amd64

all: linux_amd64

linux_amd64:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o $(CURDIR)/build/libraries/linux_amd64/libthyra.so -buildmode=c-shared internal/c/main.go