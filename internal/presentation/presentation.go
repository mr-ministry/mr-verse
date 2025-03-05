package presentation

import (
	"fmt"
	"sync"

	"github.com/mr-ministry/mr-verse/internal/bible"
)

// VersePresentation represents the current verse being presented
type VersePresentation struct {
	CurrentVerse *bible.Verse
	mu           sync.RWMutex
	observers    []func(*bible.Verse)
}

// NewVersePresentation creates a new verse presentation
func NewVersePresentation() *VersePresentation {
	return &VersePresentation{
		observers: make([]func(*bible.Verse), 0),
	}
}

// SetVerse sets the current verse and notifies all observers
func (vp *VersePresentation) SetVerse(verse *bible.Verse) {
	vp.mu.Lock()
	vp.CurrentVerse = verse
	observers := vp.observers // Copy to avoid holding lock during callbacks
	vp.mu.Unlock()

	// Notify all observers
	for _, observer := range observers {
		observer(verse)
	}
}

// GetVerse returns the current verse
func (vp *VersePresentation) GetVerse() *bible.Verse {
	vp.mu.RLock()
	defer vp.mu.RUnlock()
	return vp.CurrentVerse
}

// AddObserver adds a function to be called when the verse changes
func (vp *VersePresentation) AddObserver(observer func(*bible.Verse)) {
	vp.mu.Lock()
	defer vp.mu.Unlock()
	vp.observers = append(vp.observers, observer)
}

// FetchAndSetVerse fetches a verse from the database and sets it as the current verse
func (vp *VersePresentation) FetchAndSetVerse(translation, book string, chapter, verse int) error {
	v, err := bible.GetVerse(translation, book, chapter, verse)
	if err != nil {
		return err
	}
	vp.SetVerse(v)
	return nil
}

// FetchAndSetNextVerse fetches the next verse and sets it as the current verse
func (vp *VersePresentation) FetchAndSetNextVerse() error {
	current := vp.GetVerse()
	if current == nil {
		return fmt.Errorf("no current verse to get next from")
	}

	next, err := bible.GetNextVerse(current.Translation, current.Book, current.Chapter, current.Verse)
	if err != nil {
		return err
	}

	vp.SetVerse(next)
	return nil
}

// FetchAndSetPreviousVerse fetches the previous verse and sets it as the current verse
func (vp *VersePresentation) FetchAndSetPreviousVerse() error {
	current := vp.GetVerse()
	if current == nil {
		return fmt.Errorf("no current verse to get previous from")
	}

	prev, err := bible.GetPreviousVerse(current.Translation, current.Book, current.Chapter, current.Verse)
	if err != nil {
		return err
	}

	vp.SetVerse(prev)
	return nil
}

// SwitchTranslation switches to a different translation of the same verse
func (vp *VersePresentation) SwitchTranslation(newTranslation string) error {
	current := vp.GetVerse()
	if current == nil {
		return fmt.Errorf("no current verse to switch translation")
	}

	// Get the same verse in the new translation
	v, err := bible.GetVerse(newTranslation, current.Book, current.Chapter, current.Verse)
	if err != nil {
		return err
	}

	vp.SetVerse(v)
	return nil
}
