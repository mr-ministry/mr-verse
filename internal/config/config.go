package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// MonitorBounds represents the position and size of a monitor
type MonitorBounds struct {
	X      int
	Y      int
	Width  int
	Height int
}

// GetSecondaryMonitorBounds retrieves the secondary monitor bounds from the environment variable
// Returns nil if the environment variable is not set or is invalid
func GetSecondaryMonitorBounds() (*MonitorBounds, error) {
	boundsStr := os.Getenv("SECONDARY_MONITOR_BOUNDS")
	if boundsStr == "" {
		return nil, nil // No error, just not configured
	}

	parts := strings.Split(boundsStr, ",")
	if len(parts) != 4 {
		return nil, fmt.Errorf("invalid SECONDARY_MONITOR_BOUNDS format: expected 'x,y,width,height', got '%s'", boundsStr)
	}

	var bounds MonitorBounds
	var err error

	bounds.X, err = strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return nil, fmt.Errorf("invalid X value in SECONDARY_MONITOR_BOUNDS: %w", err)
	}

	bounds.Y, err = strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return nil, fmt.Errorf("invalid Y value in SECONDARY_MONITOR_BOUNDS: %w", err)
	}

	bounds.Width, err = strconv.Atoi(strings.TrimSpace(parts[2]))
	if err != nil {
		return nil, fmt.Errorf("invalid Width value in SECONDARY_MONITOR_BOUNDS: %w", err)
	}

	bounds.Height, err = strconv.Atoi(strings.TrimSpace(parts[3]))
	if err != nil {
		return nil, fmt.Errorf("invalid Height value in SECONDARY_MONITOR_BOUNDS: %w", err)
	}

	return &bounds, nil
}
