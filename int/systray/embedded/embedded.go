package embedded

import (
	_ "embed"
)

//go:embed logo.png
var Logo []byte

//go:embed logo_notification.png
var NotificationLogo []byte
