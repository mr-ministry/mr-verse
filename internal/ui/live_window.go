package ui

import (
	"fmt"
	"image/color"

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
	window     fyne.Window
	verseText  *widget.RichText
	reference  *widget.Label
	isOpen     bool
	app        fyne.App
	onClose    func()
}

// NewLiveWindow creates a new live window
func NewLiveWindow(app fyne.App, onClose func()) *LiveWindow {
	lw := &LiveWindow{
		app:     app,
		onClose: onClose,
		isOpen:  false,
	}
	return lw
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
	lw.window.SetOnClosed(func() {
		lw.isOpen = false
		if lw.onClose != nil {
			lw.onClose()
		}
	})

	// Create the UI components
	lw.verseText = widget.NewRichText()
	lw.verseText.Wrapping = fyne.TextWrapWord
	
	// Make the text larger for presentation
	lw.verseText.Segments = []widget.RichTextSegment{
		&widget.TextSegment{
			Style: widget.RichTextStyle{
				SizeName: theme.SizeNameHeadingText,
			},
			Text: "Welcome to Mr Verse",
		},
	}

	lw.reference = widget.NewLabel("")
	lw.reference.Alignment = fyne.TextAlignCenter
	lw.reference.TextStyle = fyne.TextStyle{Bold: true}

	// Create the layout
	content := container.NewVBox(
		container.NewCenter(lw.reference),
		widget.NewSeparator(),
		container.New(layout.NewPaddedLayout(), lw.verseText),
	)

	// Set the content
	lw.window.SetContent(content)

	// Configure for secondary monitor if specified
	bounds, err := config.GetSecondaryMonitorBounds()
	if err != nil {
		fmt.Printf("Error getting secondary monitor bounds: %v\n", err)
	}

	if bounds != nil {
		lw.window.Resize(fyne.NewSize(float32(bounds.Width), float32(bounds.Height)))
		// Position the window on the secondary monitor
		lw.window.CenterOnScreen()
		lw.window.SetFullScreen(true)
	} else {
		// Default size if no secondary monitor is specified
		lw.window.Resize(fyne.NewSize(800, 600))
		lw.window.CenterOnScreen()
	}

	lw.window.Show()
	lw.isOpen = true
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
	referenceText := fmt.Sprintf("%s %d:%d (%s)", verse.Book, verse.Chapter, verse.Verse, verse.Translation)
	lw.reference.SetText(referenceText)

	// Update the verse text
	lw.verseText.Segments = []widget.RichTextSegment{
		&widget.TextSegment{
			Style: widget.RichTextStyle{
				SizeName: theme.SizeNameHeadingText,
			},
			Text: verse.Text,
		},
	}
	lw.verseText.Refresh()
}

// SetBackground sets the background color of the live window
func (lw *LiveWindow) SetBackground(color color.Color) {
	if !lw.isOpen {
		return
	}

	bg := canvas.NewRectangle(color)
	lw.window.SetContent(container.NewMax(bg, lw.window.Content()))
}
