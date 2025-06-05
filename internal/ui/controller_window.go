package ui

import (
	"fmt"
	"log"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/mr-ministry/mr-verse/internal/bible"
	"github.com/mr-ministry/mr-verse/internal/config"
	"github.com/mr-ministry/mr-verse/internal/presentation"
)

// ControllerWindow represents the main control window
type ControllerWindow struct {
	window            fyne.Window
	app               fyne.App
	liveWindow        *LiveWindow
	versePresentation *presentation.VersePresentation
	searchEntry       *widget.Entry
	translationSelect *widget.Select
	statusLabel       *widget.Label
	currentVerseLabel *widget.Label
}

// RunApp initializes and runs the application
func RunApp() {
	// Initialize the app
	a := app.New()
	w := a.NewWindow("Mr Verse - Controller")
	w.Resize(fyne.NewSize(800, 600))

	// Initialize and seed the database
	if err := initializeDatabases(w); err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Create the controller window
	controller := &ControllerWindow{
		window:            w,
		app:               a,
		versePresentation: presentation.NewVersePresentation(),
	}

	// Create the live window
	controller.liveWindow = NewLiveWindow(a, func() {
		controller.updateLiveWindowStatus(false)
	})

	// Set up the UI
	controller.setupUI()

	// Show the window
	w.ShowAndRun()

	// Clean up
	bible.CloseDB()
}

// initializeDatabases initializes and seeds the Bible database
func initializeDatabases(w fyne.Window) error {
	// Initialize the database
	if err := bible.InitDB(); err != nil {
		dialog.ShowError(fmt.Errorf("failed to initialize database: %w", err), w)
		return err
	}

	// Seed the database with Bible data
	if err := bible.SeedBibleData(); err != nil {
		dialog.ShowError(fmt.Errorf("failed to seed Bible data: %w", err), w)
		log.Printf("Failed to seed Bible data: %v", err)
		// Not a fatal error, can continue
	}

	return nil
}

// setupUI sets up the user interface
// TODO: Default verse to show should be John 3:16
func (c *ControllerWindow) setupUI() {
	// Create the search entry
	c.searchEntry = widget.NewEntry()
	// c.searchEntry.SetPlaceHolder("Enter Bible reference (e.g., John 3:16)")
	c.searchEntry.SetText("Esther 8:9") // use this to set the default verse

	// Create the search button
	searchButton := widget.NewButton("Search", func() {
		c.searchVerse()
	})
	c.searchEntry.OnSubmitted = func(s string) {
		c.searchVerse()
	}

	// Create the translation select
	c.translationSelect = widget.NewSelect([]string{"Loading..."}, func(s string) {
		c.switchTranslation(s)
	})

	// Load available translations
	go c.loadTranslations()

	// Create the navigation buttons
	prevButton := widget.NewButton("Previous Verse", func() {
		c.navigateToPreviousVerse()
	})
	nextButton := widget.NewButton("Next Verse", func() {
		c.navigateToNextVerse()
	})

	// Create the live window control button
	// TODO: Change the button color when the live window is open
	liveWindowButton := widget.NewButton("Go Live", func() {
		if c.liveWindow.IsOpen() {
			c.liveWindow.Close()
		} else {
			c.liveWindow.Open()
			c.updateLiveWindowStatus(true)
		}
	})

	// Create the update live window button
	updateLiveButton := widget.NewButton("Update Live Window", func() {
		c.updateLiveWindow()
	})

	// Create the settings button
	settingsButton := widget.NewButton("Settings", func() {
		c.showSettingsDialog()
	})

	// Create the status label
	c.statusLabel = widget.NewLabel("Offline")
	c.currentVerseLabel = widget.NewLabel("No verse selected")

	// Create the layout
	searchContainer := container.NewBorder(nil, nil, nil, searchButton, c.searchEntry)

	controlsContainer := container.NewVBox(
		widget.NewLabel("Bible Translation:"),
		c.translationSelect,
		container.NewHBox(prevButton, nextButton),
		container.NewHBox(liveWindowButton, updateLiveButton),
		settingsButton,
	)

	statusContainer := container.NewHBox(
		widget.NewLabel("Status:"),
		c.statusLabel,
		widget.NewLabel("Current Verse:"),
		c.currentVerseLabel,
	)

	// Main layout
	mainContainer := container.NewBorder(
		searchContainer,
		statusContainer,
		nil,
		nil,
		container.New(layout.NewCenterLayout(), controlsContainer),
	)

	c.window.SetContent(mainContainer)

	// Register as an observer for verse changes
	c.versePresentation.AddObserver(func(verse *bible.Verse) {
		if verse != nil {
			c.updateCurrentVerseLabel(verse)

			// Update the live window if it's open
			if c.liveWindow.IsOpen() {
				c.liveWindow.UpdateVerse(verse)
			}
		}
	})
}

// loadTranslations loads the available Bible translations
func (c *ControllerWindow) loadTranslations() {
	translations, err := bible.GetAvailableTranslations()
	if err != nil {
		// Use a goroutine to show the error dialog on the main thread
		go func() {
			dialog.ShowError(
				fmt.Errorf("failed to load translations: %w", err),
				c.window,
			)
		}()
		return
	}

	// Use a goroutine to update the UI on the main thread
	go func() {
		if len(translations) > 0 {
			c.translationSelect.Options = translations
			c.translationSelect.SetSelected(translations[0])
		} else {
			c.translationSelect.Options = []string{"No translations available"}
		}
		c.translationSelect.Refresh()
	}()
}

// searchVerse searches for a Bible verse
// TODO: Allow for fuzzy search using keywords from verse reference or verse text
func (c *ControllerWindow) searchVerse() {
	reference := c.searchEntry.Text
	if reference == "" {
		dialog.ShowInformation("Error", "Please enter a Bible reference", c.window)
		return
	}

	// Parse the reference
	book, chapter, verse, err := bible.ParseBibleReference(reference)
	if err != nil {
		dialog.ShowError(fmt.Errorf("invalid Bible reference: %w", err), c.window)
		return
	}

	// Get the selected translation
	translation := c.translationSelect.Selected
	if translation == "" || translation == "Loading..." ||
		translation == "No translations available" {
		dialog.ShowInformation("Error", "Please select a valid translation", c.window)
		return
	}

	// Fetch and set the verse
	err = c.versePresentation.FetchAndSetVerse(translation, book, chapter, verse)
	if err != nil {
		dialog.ShowError(fmt.Errorf("failed to fetch verse: %w", err), c.window)
		return
	}
}

// navigateToNextVerse navigates to the next verse
func (c *ControllerWindow) navigateToNextVerse() {
	if c.versePresentation.GetVerse() == nil {
		dialog.ShowInformation("Error", "No current verse selected", c.window)
		return
	}

	err := c.versePresentation.FetchAndSetNextVerse()
	if err != nil {
		dialog.ShowError(fmt.Errorf("failed to fetch next verse: %w", err), c.window)
		return
	}
}

// navigateToPreviousVerse navigates to the previous verse
func (c *ControllerWindow) navigateToPreviousVerse() {
	if c.versePresentation.GetVerse() == nil {
		dialog.ShowInformation("Error", "No current verse selected", c.window)
		return
	}

	err := c.versePresentation.FetchAndSetPreviousVerse()
	if err != nil {
		dialog.ShowError(
			fmt.Errorf("failed to fetch previous verse: %w", err),
			c.window,
		)
		return
	}
}

// switchTranslation switches to a different translation
func (c *ControllerWindow) switchTranslation(translation string) {
	if c.versePresentation.GetVerse() == nil {
		// No verse selected yet, nothing to do
		return
	}

	err := c.versePresentation.SwitchTranslation(translation)
	if err != nil {
		dialog.ShowError(fmt.Errorf("failed to switch translation: %w", err), c.window)
		return
	}
}

// updateLiveWindow updates the live window with the current verse
func (c *ControllerWindow) updateLiveWindow() {
	if !c.liveWindow.IsOpen() {
		dialog.ShowInformation("Error", "Live window is not open", c.window)
		return
	}

	verse := c.versePresentation.GetVerse()
	if verse == nil {
		dialog.ShowInformation("Error", "No verse selected", c.window)
		return
	}

	c.liveWindow.UpdateVerse(verse)
}

// updateLiveWindowStatus updates the status label based on the live window state
// TODO: Set text colors depending on status
func (c *ControllerWindow) updateLiveWindowStatus(isOpen bool) {
	if isOpen {
		c.statusLabel.SetText("Live")
	} else {
		c.statusLabel.SetText("Offline")
	}
}

// updateCurrentVerseLabel updates the current verse label
func (c *ControllerWindow) updateCurrentVerseLabel(verse *bible.Verse) {
	if verse == nil {
		c.currentVerseLabel.SetText("No verse selected")
		return
	}

	c.currentVerseLabel.SetText(
		fmt.Sprintf(
			"%s %d:%d (%s)",
			verse.Book,
			verse.Chapter,
			verse.Verse,
			verse.Translation,
		),
	)
}

// showSettingsDialog displays the settings dialog
func (c *ControllerWindow) showSettingsDialog() {
	// Get current bounds from preferences
	currentBounds := config.GetMonitorBounds(c.app.Preferences())

	// Default values
	var defaultX, defaultY, defaultWidth, defaultHeight int
	if currentBounds != nil {
		defaultX = currentBounds.X
		defaultY = currentBounds.Y
		defaultWidth = currentBounds.Width
		defaultHeight = currentBounds.Height
	} else {
		// Default values if not set
		defaultX = 1920 // Common values for a secondary monitor
		defaultY = 0
		defaultWidth = 1920
		defaultHeight = 1200 // Standard 16:10 resolution
		// defaultHeight = 1080 // Standard 16:9 resolution
	}

	// Create form fields
	xEntry := widget.NewEntry()
	xEntry.SetText(strconv.Itoa(defaultX))

	yEntry := widget.NewEntry()
	yEntry.SetText(strconv.Itoa(defaultY))

	widthEntry := widget.NewEntry()
	widthEntry.SetText(strconv.Itoa(defaultWidth))

	heightEntry := widget.NewEntry()
	heightEntry.SetText(strconv.Itoa(defaultHeight))

	// Create the form
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "X Position", Widget: xEntry},
			{Text: "Y Position", Widget: yEntry},
			{Text: "Width", Widget: widthEntry},
			{Text: "Height", Widget: heightEntry},
		},
		OnSubmit: func() {
			// Parse the values
			x, err1 := strconv.Atoi(xEntry.Text)
			y, err2 := strconv.Atoi(yEntry.Text)
			width, err3 := strconv.Atoi(widthEntry.Text)
			height, err4 := strconv.Atoi(heightEntry.Text)

			// Check for errors
			if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
				dialog.ShowError(
					fmt.Errorf("invalid values, please enter numbers only"),
					c.window,
				)
				return
			}

			// Save to preferences
			bounds := &config.MonitorBounds{
				X:      x,
				Y:      y,
				Width:  width,
				Height: height,
			}
			config.SaveMonitorBounds(c.app.Preferences(), bounds)

			// If live window is open, refresh it with new settings
			if c.liveWindow.IsOpen() {
				// Update theme with new size before closing/reopening
				if currentTheme, ok := c.app.Settings().Theme().(*presentationTheme); ok {
					currentTheme.UpdateWindowSize(
						fyne.NewSize(float32(width), float32(height)),
					)
				}
				c.liveWindow.Close()
				c.liveWindow.Open()
			}

			dialog.ShowInformation(
				"Settings Saved",
				"Monitor settings have been saved.",
				c.window,
			)
		},
		OnCancel: func() {
			// Do nothing
		},
		SubmitText: "Save",
		CancelText: "Cancel",
	}

	// Create and show the dialog
	dialog.ShowCustom("Secondary Monitor Settings", "Close", form, c.window)
}
