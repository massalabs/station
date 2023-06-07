package utils

import (
	"log"
	"net/url"

	"fyne.io/fyne/v2"
)

func OpenURL(app *fyne.App, urlToOpen string) {
	u, err := url.Parse(urlToOpen)
	if err != nil {
		log.Println("Error parsing URL:", err)
	}

	err = (*app).OpenURL(u)
	if err != nil {
		log.Println("Error opening URL:", err)
	}
}
