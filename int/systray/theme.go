package systray

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"github.com/massalabs/station/int/systray/embedded"
)

// Colors of the Station web UI default theme (react-ui-kit massa-station-preset, dark-v2).
//
//nolint:gochecknoglobals,mnd // Palette constants.
var (
	colorBackground     = color.NRGBA{R: 0x19, G: 0x18, B: 0x29, A: 0xff} // primary
	colorSurface        = color.NRGBA{R: 0x2e, G: 0x37, B: 0x43, A: 0xff} // secondary
	colorSurfacePressed = color.NRGBA{R: 0x1a, G: 0x20, B: 0x22, A: 0xff} // tertiary
	colorBrand          = color.NRGBA{R: 0x1a, G: 0xe1, B: 0x9d, A: 0xff} // s-success
	colorText           = color.NRGBA{R: 0xda, G: 0xda, B: 0xda, A: 0xff} // f-primary
	colorTextSecondary  = color.NRGBA{R: 0xf1, G: 0xf1, B: 0xf1, A: 0xff} // f-secondary
	colorButton         = color.NRGBA{R: 0x00, G: 0x43, B: 0xff, A: 0xff} // c-default
	colorDisabled       = color.NRGBA{R: 0x6b, G: 0x72, B: 0x80, A: 0xff} // c-disabled-1
	colorDisabledButton = color.NRGBA{R: 0x37, G: 0x41, B: 0x51, A: 0xff} // c-disabled-2
	colorError          = color.NRGBA{R: 0xff, G: 0x4e, B: 0x4e, A: 0xff} // s-error
	colorWarning        = color.NRGBA{R: 0xff, G: 0xa5, B: 0x1c, A: 0xff} // s-warning
	colorInfo           = color.NRGBA{R: 0x4a, G: 0xb2, B: 0xff, A: 0xff} // s-info-1
	colorHover          = color.NRGBA{R: 0xda, G: 0xda, B: 0xda, A: 0x26} // neutral/15, as the web secondary buttons
	colorShadow         = color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x66}
)

// stationTheme is the Fyne theme matching the Station web UI, used by the systray windows.
// It always uses the dark palette, as the web UI does by default, regardless of the system variant.
type stationTheme struct {
	fonts map[fyne.TextStyle]fyne.Resource
}

func newStationTheme() fyne.Theme {
	return &stationTheme{
		fonts: map[fyne.TextStyle]fyne.Resource{
			{}:                         fyne.NewStaticResource("Poppins-Regular.ttf", embedded.FontRegular),
			{Bold: true}:               fyne.NewStaticResource("Poppins-SemiBold.ttf", embedded.FontBold),
			{Italic: true}:             fyne.NewStaticResource("Poppins-Italic.ttf", embedded.FontItalic),
			{Bold: true, Italic: true}: fyne.NewStaticResource("Poppins-SemiBoldItalic.ttf", embedded.FontBoldItalic),
		},
	}
}

//nolint:cyclop // Mapping of the palette.
func (t *stationTheme) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground, theme.ColorNameOverlayBackground, theme.ColorNameMenuBackground,
		theme.ColorNameHeaderBackground:
		return colorBackground
	case theme.ColorNameInputBackground, theme.ColorNameButton, theme.ColorNameScrollBarBackground:
		return colorSurface
	case theme.ColorNameForeground, theme.ColorNamePlaceHolder:
		return colorText
	case theme.ColorNameForegroundOnPrimary, theme.ColorNameForegroundOnError,
		theme.ColorNameForegroundOnSuccess, theme.ColorNameForegroundOnWarning:
		return colorTextSecondary
	case theme.ColorNamePrimary, theme.ColorNameFocus:
		return colorButton
	case theme.ColorNameHyperlink:
		return colorInfo
	case theme.ColorNameSuccess:
		return colorBrand
	case theme.ColorNameError:
		return colorError
	case theme.ColorNameWarning:
		return colorWarning
	case theme.ColorNameDisabled:
		return colorDisabled
	case theme.ColorNameDisabledButton:
		return colorDisabledButton
	case theme.ColorNameHover, theme.ColorNameSelection:
		return colorHover
	case theme.ColorNamePressed:
		return colorSurfacePressed
	case theme.ColorNameInputBorder, theme.ColorNameSeparator, theme.ColorNameScrollBar:
		return colorSurface
	case theme.ColorNameShadow:
		return colorShadow
	default:
		return theme.DefaultTheme().Color(name, theme.VariantDark)
	}
}

func (t *stationTheme) Font(style fyne.TextStyle) fyne.Resource {
	if style.Monospace || style.Symbol {
		return theme.DefaultTheme().Font(style)
	}

	return t.fonts[fyne.TextStyle{Bold: style.Bold, Italic: style.Italic}]
}

func (t *stationTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

//nolint:mnd // Sizes of the Station web UI.
func (t *stationTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNameText:
		return 14
	case theme.SizeNameSubHeadingText:
		return 18
	case theme.SizeNameHeadingText:
		return 22
	case theme.SizeNameInputRadius, theme.SizeNameSelectionRadius:
		// Buttons are rounded-lg.
		return 8
	case theme.SizeNameInnerPadding:
		// Buttons are 48px high.
		return 12
	case theme.SizeNamePadding:
		return 6
	default:
		return theme.DefaultTheme().Size(name)
	}
}
