package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type presentationTheme struct{}

var _ fyne.Theme = (*presentationTheme)(nil)

// NewPresentationTheme creates a new theme instance for the presentation window
func NewPresentationTheme() fyne.Theme {
	return &presentationTheme{}
}

// Color returns the color for the specified name and theme variant
func (t *presentationTheme) Color(
	name fyne.ThemeColorName,
	variant fyne.ThemeVariant,
) color.Color {
	if name == theme.ColorNameForeground {
		return color.White // Force white text
	}
	return theme.DefaultTheme().Color(name, variant)
}

// Font returns the font resource for the specified text style
func (t *presentationTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

// Icon returns the icon resource for the specified icon name
func (t *presentationTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

// Size returns the size for the specified size name
// TODO: Make the max font size dynamic based on the window size and/or screen resolution
func (t *presentationTheme) Size(name fyne.ThemeSizeName) float32 {
	if name == theme.SizeNameSubHeadingText {
		return 30 // Much larger text size for presentation
	}
	if name == theme.SizeNameHeadingText {
		return 60 // Even larger for headings
	}
	return theme.DefaultTheme().Size(name)
}
