package ui

import (
	"image/color"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// presentationTheme customizes the appearance of the presentation window
type presentationTheme struct {
	windowSize fyne.Size
}

var _ fyne.Theme = (*presentationTheme)(nil)

// NewPresentationTheme creates a new theme instance for the presentation window
func NewPresentationTheme() fyne.Theme {
	return &presentationTheme{
		windowSize: fyne.NewSize(1920, 1600), // Default size
	}
}

// NewPresentationThemeWithSize creates a new theme instance with the specified window size
func NewPresentationThemeWithSize(size fyne.Size) fyne.Theme {
	return &presentationTheme{
		windowSize: size,
	}
}

// Color returns white for text foreground, otherwise defaults
func (t *presentationTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == theme.ColorNameForeground {
		return color.White // Force white text
	}
	return theme.DefaultTheme().Color(name, variant)
}

// Font delegates to the default theme
func (t *presentationTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

// Icon delegates to the default theme
func (t *presentationTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

// Size returns dynamically calculated text sizes based on screen resolution
func (t *presentationTheme) Size(name fyne.ThemeSizeName) float32 {
	if name != theme.SizeNameHeadingText && name != theme.SizeNameSubHeadingText {
		return theme.DefaultTheme().Size(name)
	}
	
	// Base sizes calibrated for 1920x1080 resolution
	baseHeadingSize := float32(60)
	baseSubHeadingSize := float32(30)
	
	// Calculate the scale factor - use sqrt of area ratio for balanced scaling
	referenceArea := float32(1920 * 1080)
	actualArea := t.windowSize.Width * t.windowSize.Height
	
	// Calculate scale factor with boundaries to prevent extreme sizes
	scale := float32(math.Sqrt(float64(actualArea / referenceArea)))
	scale = float32(math.Max(0.5, math.Min(float64(scale), 1.5))) // Limit scale between 0.5 and 1.5
	
	if name == theme.SizeNameHeadingText {
		return baseHeadingSize * scale
	}
	return baseSubHeadingSize * scale
}

// UpdateWindowSize allows updating the window size after theme creation
func (t *presentationTheme) UpdateWindowSize(size fyne.Size) {
	t.windowSize = size
}
