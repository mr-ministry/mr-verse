package ui

import (
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/mr-ministry/mr-verse/internal/bible"
	"github.com/mr-ministry/mr-verse/internal/config"
)

// LiveWindow represents the presentation window
type LiveWindow struct {
	window    fyne.Window
	app       fyne.App
	verseText *widget.RichText
	reference *widget.RichText
	isOpen    bool
	onClose   func()
}

// NewLiveWindow creates a new live window
func NewLiveWindow(app fyne.App, onClose func()) *LiveWindow {
	return &LiveWindow{
		app:     app,
		onClose: onClose,
		isOpen:  false,
	}
}

// IsOpen returns whether the live window is open
func (lw *LiveWindow) IsOpen() bool {
	return lw.isOpen
}

// Open opens the live window
func (lw *LiveWindow) Open() {
	if lw.isOpen {
		lw.window.Show()
		return
	}

	// Create the window
	lw.window = lw.app.NewWindow("Mr Verse - Live Presentation")
	
	// Get monitor bounds from preferences
	bounds := config.GetMonitorBounds(lw.app.Preferences())
	
	// Calculate size for the window
	var windowSize fyne.Size
	if bounds != nil {
		windowSize = fyne.NewSize(float32(bounds.Width), float32(bounds.Height))
	} else {
		// Default size if no secondary monitor is specified
		windowSize = fyne.NewSize(800, 600)
	}
	
	// Set our custom theme for larger text with size awareness
	customTheme := NewPresentationThemeWithSize(windowSize)
	lw.app.Settings().SetTheme(customTheme)
	
	lw.window.SetOnClosed(func() {
		lw.isOpen = false
		if lw.onClose != nil {
			lw.onClose()
		}
	})

	// Create the UI components
	lw.setupUI()

	// Set the content
	lw.window.Resize(windowSize)
	lw.window.CenterOnScreen()
	
	if bounds != nil {
		lw.window.SetFullScreen(true)
	}
	
	// Set up a listener for size changes and update the theme dynamically
	var lastSize fyne.Size = windowSize
	go lw.monitorWindowSize(lastSize)

	lw.window.Show()
	lw.isOpen = true
}

// setupUI creates the UI components for the live window
func (lw *LiveWindow) setupUI() {
	lw.verseText = widget.NewRichText()
	lw.verseText.Wrapping = fyne.TextWrapWord

	// Make the text larger for presentation with white color
	lw.verseText.Segments = []widget.RichTextSegment{
		&widget.TextSegment{
			Style: widget.RichTextStyle{
				TextStyle: fyne.TextStyle{
					Bold: true,
				},
				Alignment: fyne.TextAlignCenter,
				SizeName:  theme.SizeNameHeadingText,
			},
			Text: "JESUS IS KING",
		},
	}

	lw.reference = widget.NewRichText()
	lw.reference.Segments = []widget.RichTextSegment{
		&widget.TextSegment{
			Style: widget.RichTextStyle{
				TextStyle: fyne.TextStyle{
					Bold: true,
				},
				SizeName:  theme.SizeNameSubHeadingText,
				Alignment: fyne.TextAlignCenter,
			},
			Text: "...",
		},
	}

	// Create the layout
	content := container.NewVBox(
		container.NewCenter(lw.reference),
		widget.NewSeparator(),
		layout.NewSpacer(),
		container.New(layout.NewPaddedLayout(), lw.verseText),
		layout.NewSpacer(),
	)

	// Set dark background
	bg := canvas.NewRectangle(color.Black)
	mainContent := container.NewStack(bg, content)

	// Set the content
	lw.window.SetContent(mainContent)
}

// monitorWindowSize monitors window size changes and updates the theme accordingly
func (lw *LiveWindow) monitorWindowSize(initialSize fyne.Size) {
	lastSize := initialSize
	for lw.isOpen {
		// Check size every 500ms
		time.Sleep(500 * time.Millisecond)
		if !lw.isOpen {
			break
		}
		
		currentSize := lw.window.Canvas().Size()
		if currentSize.Width != lastSize.Width || currentSize.Height != lastSize.Height {
			// Size has changed, update the theme
			if currentTheme, ok := lw.app.Settings().Theme().(*presentationTheme); ok {
				currentTheme.UpdateWindowSize(currentSize)
				// Force refresh text
				lw.verseText.Refresh()
				lw.reference.Refresh()
			}
			lastSize = currentSize
		}
	}
}

// Close closes the live window
func (lw *LiveWindow) Close() {
	if lw.isOpen && lw.window != nil {
		lw.window.Close()
		lw.isOpen = false
	}
}

// UpdateVerse updates the verse displayed in the live window
func (lw *LiveWindow) UpdateVerse(verse *bible.Verse) {
	if !lw.isOpen || verse == nil {
		return
	}

	// Update the reference
	referenceText := fmt.Sprintf(
		"%s %d:%d (%s)",
		verse.Book,
		verse.Chapter,
		verse.Verse,
		verse.Translation,
	)
	lw.reference.Segments = []widget.RichTextSegment{
		&widget.TextSegment{
			Style: widget.RichTextStyle{
				TextStyle: fyne.TextStyle{
					Bold: true,
				},
				SizeName:  theme.SizeNameSubHeadingText,
				Alignment: fyne.TextAlignCenter,
			},
			Text: string(referenceText),
		},
	}

	// Update the verse text with white color
	lw.verseText.Segments = []widget.RichTextSegment{
		&widget.TextSegment{
			Style: widget.RichTextStyle{
				TextStyle: fyne.TextStyle{
					Bold: true,
				},
				Alignment: fyne.TextAlignCenter,
				SizeName:  theme.SizeNameHeadingText,
			},
			Text: verse.Text,
		},
	}
	lw.verseText.Wrapping = fyne.TextWrapWord

	lw.verseText.Refresh()
	lw.reference.Refresh()
}

// SetBackground sets the background color of the live window
func (lw *LiveWindow) SetBackground(color color.Color) {
	if !lw.isOpen {
		return
	}

	bg := canvas.NewRectangle(color)
	lw.window.SetContent(container.NewStack(bg, lw.window.Content()))
}
