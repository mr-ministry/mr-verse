package config

import (
	"fyne.io/fyne/v2"
)

// MonitorBounds represents the position and size of a monitor
type MonitorBounds struct {
	X      int
	Y      int
	Width  int
	Height int
}

// Preference keys
const (
	PrefKeyMonitorX      = "secondaryMonitor.x"
	PrefKeyMonitorY      = "secondaryMonitor.y"
	PrefKeyMonitorWidth  = "secondaryMonitor.width"
	PrefKeyMonitorHeight = "secondaryMonitor.height"
)

// GetMonitorBounds retrieves the secondary monitor bounds from app preferences
func GetMonitorBounds(preferences fyne.Preferences) *MonitorBounds {
	// Check if all required values are present
	x := preferences.Int(PrefKeyMonitorX)
	y := preferences.Int(PrefKeyMonitorY)
	width := preferences.Int(PrefKeyMonitorWidth)
	height := preferences.Int(PrefKeyMonitorHeight)

	// If any value is 0, consider it not configured
	// (since 0 is the default for int preferences that aren't set)
	if width == 0 || height == 0 {
		return nil
	}

	// Get values from preferences
	return &MonitorBounds{
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
	}
}

// SaveMonitorBounds saves the monitor bounds to app preferences
func SaveMonitorBounds(preferences fyne.Preferences, bounds *MonitorBounds) {
	if bounds == nil {
		return
	}

	preferences.SetInt(PrefKeyMonitorX, bounds.X)
	preferences.SetInt(PrefKeyMonitorY, bounds.Y)
	preferences.SetInt(PrefKeyMonitorWidth, bounds.Width)
	preferences.SetInt(PrefKeyMonitorHeight, bounds.Height)
}
