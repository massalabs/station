package main

import (
	"fyne.io/fyne/v2/app"
	"github.com/massalabs/thyra/int/api"
)

func main() {
	myApp := app.New()

	go api.StartServer(&myApp)

	myApp.Run()
}
