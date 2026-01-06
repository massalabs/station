package main

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"github.com/massalabs/station/int/config"
	"github.com/massalabs/station/int/systray/embedded"
)

func init() {
	app.SetMetadata(fyne.AppMetadata{
		ID:      "net.massalabs.massastation",
		Name:    "MassaStation",
		Version: config.Version,
		Build:   1,
		Icon: &fyne.StaticResource{
			StaticName:    "logo.png",
			StaticContent: embedded.Logo,
		},
		Release: !strings.Contains(config.Version, "dev"),
	})
}
