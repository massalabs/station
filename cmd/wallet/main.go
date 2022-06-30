package main

import (
	"fmt"

	"github.com/massalabs/thyra/pkg/wallet"
)

func main() {
	w, err := wallet.New("testing")
	if err != nil {
		panic(err)
	}

	fmt.Println("init", w.KeyPairs[0].PrivateKey)

	// this panic because serialization of unprotected keypair is forbidden
	/*yaml, err := w.YAML()
	if err != nil {
		panic(err)
	}

	fmt.Println("yaml", string(yaml))*/

	err = w.Protect("MyAwesomePassword", 0)
	if err != nil {
		panic(err)
	}

	fmt.Println("protected", w.KeyPairs[0].PrivateKey)

	yaml, err := w.YAML()
	if err != nil {
		panic(err)
	}

	fmt.Println("yaml", string(yaml))

	err = w.Unprotect("MyAwesomePassword", 0)
	if err != nil {
		panic(err)
	}

	fmt.Println("unprotected", w.KeyPairs[0].PrivateKey)

	err = w.Protect("MyAwesomePassword", 0)
	if err != nil {
		panic(err)
	}

	err = w.Unprotect("WrongPassword", 0)
	if err != nil {
		panic(err)
	}

	fmt.Println("wrong password", w.KeyPairs[0].PrivateKey)

	w2, err := wallet.FromYAML(yaml)
	if err != nil {
		panic(err)
	}

	fmt.Println("unserialized", w2.KeyPairs[0].PrivateKey)

	err = w2.Unprotect("MyAwesomePassword", 0)
	if err != nil {
		panic(err)
	}

	fmt.Println("unserialized and unprotected", w2.KeyPairs[0].PrivateKey)
}
