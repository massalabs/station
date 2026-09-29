package embedded

import (
	_ "embed"
)

//go:embed logo.png
var Logo []byte

//go:embed logo_notification.png
var NotificationLogo []byte

// Poppins fonts, used by the Station web UI. Licensed under the SIL Open Font License (see fonts/Poppins-OFL.txt).

//go:embed fonts/Poppins-Regular.ttf
var FontRegular []byte

//go:embed fonts/Poppins-SemiBold.ttf
var FontBold []byte

//go:embed fonts/Poppins-Italic.ttf
var FontItalic []byte

//go:embed fonts/Poppins-SemiBoldItalic.ttf
var FontBoldItalic []byte
