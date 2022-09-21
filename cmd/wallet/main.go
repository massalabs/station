package main

import (
	"fmt"

	"github.com/massalabs/thyra/pkg/wallet"
)

func main() {
	wlt, err := wallet.New("testing")
	if err != nil {
		panic(err)
	}

	//nolint:forbidigo
	fmt.Println("init", wlt.KeyPairs[0].PrivateKey)

	err = wlt.Protect("MyAwesomePassword", 0)
	if err != nil {
		panic(err)
	}

	//nolint:forbidigo
	fmt.Println("protected", wlt.KeyPairs[0].PrivateKey)

	yaml, err := wlt.YAML()
	if err != nil {
		panic(err)
	}

	//nolint:forbidigo
	fmt.Println("yaml", string(yaml))

	err = wlt.Unprotect("MyAwesomePassword", 0)
	if err != nil {
		panic(err)
	}

	//nolint:forbidigo
	fmt.Println("unprotected", wlt.KeyPairs[0].PrivateKey)

	err = wlt.Protect("MyAwesomePassword", 0)
	if err != nil {
		panic(err)
	}

	err = wlt.Unprotect("WrongPassword", 0)
	if err == nil {
		panic("using wrong password shall be detected")
	}

	wlt2, err := wallet.FromYAML(yaml)
	if err != nil {
		panic(err)
	}

	//nolint:forbidigo
	fmt.Println("unserialized", wlt2.KeyPairs[0].PrivateKey)

	err = wlt2.Unprotect("MyAwesomePassword", 0)
	if err != nil {
		panic(err)
	}

	//nolint:forbidigo
	fmt.Println("unserialized and unprotected", wlt2.KeyPairs[0].PrivateKey)
}
