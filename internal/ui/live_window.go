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
	window        fyne.Window
	app           fyne.App
	verseText     *widget.RichText
	reference     *widget.RichText
	isOpen        bool
	onClose       func()
	verseTextSize float32
}

// NewLiveWindow creates a new live window
func NewLiveWindow(app fyne.App, onClose func()) *LiveWindow {
	lw := &LiveWindow{
		app:           app,
		onClose:       onClose,
		isOpen:        false,
		verseTextSize: 36, // Default verse text size
	}
	return lw
}

// SetVerseTextSize sets the text size for verses
func (lw *LiveWindow) SetVerseTextSize(size float32) {
	lw.verseTextSize = size
	if lw.isOpen && lw.verseText != nil {
		lw.updateVerseTextStyle()
		lw.verseText.Refresh()
	}
}

// updateVerseTextStyle updates the style of the verse text
func (lw *LiveWindow) updateVerseTextStyle() {
	if lw.verseText == nil || len(lw.verseText.Segments) == 0 {
		return
	}

	if textSegment, ok := lw.verseText.Segments[0].(*widget.TextSegment); ok {
		textSegment.Style.SizeName = theme.SizeNameHeadingText
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
	// Set our custom theme for larger text
	lw.app.Settings().SetTheme(NewPresentationTheme())
	lw.window.SetOnClosed(func() {
		lw.isOpen = false
		if lw.onClose != nil {
			lw.onClose()
		}
	})

	// Create the UI components
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
				SizeName:  theme.SizeNameHeadingText, // This will use our custom theme's large size
			},
			Text: "Welcome to Mr-Verse",
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
			Text: "@mrjxtr",
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
