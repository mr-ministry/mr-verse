package bible

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Verse represents a single Bible verse
type Verse struct {
	ID          int    `json:"id,omitempty"`
	Translation string `json:"translation"`
	Book        string `json:"book"`
	Chapter     int    `json:"chapter"`
	Verse       int    `json:"verse"`
	Text        string `json:"text"`
}

// BibleData represents the structure of the Bible JSON files
type BibleData struct {
	Version string                            `json:"version"`
	Books   map[string]map[string]ChapterData `json:"books"`
}

// ChapterData represents the structure of a chapter in the Bible JSON files
type ChapterData struct {
	Header string            `json:"header"`
	Verses map[string]string `json:"verses"`
}

// GetChapterHeader fetches a localized chapter header for a translation/book/chapter.
// Returns ok=false when there is no stored header.
func GetChapterHeader(translation, book string, chapter int) (header string, ok bool, err error) {
	query := `
		SELECT header
		FROM chapter_headers
		WHERE translation = ? AND book = ? AND chapter = ?
		LIMIT 1
	`
	row := DB.QueryRow(query, translation, book, chapter)

	var h string
	if scanErr := row.Scan(&h); scanErr != nil {
		if scanErr == sql.ErrNoRows {
			return "", false, nil
		}
		return "", false, scanErr
	}
	return h, true, nil
}

// GetAvailableTranslations returns a list of available Bible translations
func GetAvailableTranslations() ([]string, error) {
	rows, err := DB.Query("SELECT DISTINCT translation FROM bible ORDER BY translation")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var translations []string
	for rows.Next() {
		var translation string
		if err := rows.Scan(&translation); err != nil {
			return nil, err
		}
		translations = append(translations, translation)
	}

	return translations, nil
}

// GetVerse retrieves a specific verse from the database
func GetVerse(translation, book string, chapter, verse int) (*Verse, error) {
	query := `
		SELECT id, translation, book, chapter, verse, text
		FROM bible
		WHERE translation = ? AND book = ? AND chapter = ? AND verse = ?
	`
	row := DB.QueryRow(query, translation, book, chapter, verse)

	var v Verse
	err := row.Scan(&v.ID, &v.Translation, &v.Book, &v.Chapter, &v.Verse, &v.Text)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf(
				"verse not found: %s %s %d:%d",
				translation,
				book,
				chapter,
				verse,
			)
		}
		return nil, err
	}

	return &v, nil
}

// GetNextVerse retrieves the next verse in sequence
func GetNextVerse(translation, book string, chapter, verse int) (*Verse, error) {
	// First try to get the next verse in the same chapter
	query := `
		SELECT id, translation, book, chapter, verse, text
		FROM bible
		WHERE translation = ? AND book = ? AND chapter = ? AND verse > ?
		ORDER BY verse ASC
		LIMIT 1
	`
	row := DB.QueryRow(query, translation, book, chapter, verse)

	var v Verse
	err := row.Scan(&v.ID, &v.Translation, &v.Book, &v.Chapter, &v.Verse, &v.Text)
	if err == nil {
		return &v, nil
	}

	if err != sql.ErrNoRows {
		return nil, err
	}

	// If no next verse in the same chapter, try to get the first verse of the next chapter
	query = `
		SELECT id, translation, book, chapter, verse, text
		FROM bible
		WHERE translation = ? AND book = ? AND chapter > ?
		ORDER BY chapter ASC, verse ASC
		LIMIT 1
	`
	row = DB.QueryRow(query, translation, book, chapter)
	err = row.Scan(&v.ID, &v.Translation, &v.Book, &v.Chapter, &v.Verse, &v.Text)
	if err == nil {
		return &v, nil
	}

	if err != sql.ErrNoRows {
		return nil, err
	}

	// If no next chapter in the same book, try to get the first verse of the next book
	query = `
		SELECT id, translation, book, chapter, verse, text
		FROM bible
		WHERE translation = ? AND book > ?
		ORDER BY book ASC, chapter ASC, verse ASC
		LIMIT 1
	`
	row = DB.QueryRow(query, translation, book)
	err = row.Scan(&v.ID, &v.Translation, &v.Book, &v.Chapter, &v.Verse, &v.Text)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf(
				"no next verse found: %s %s %d:%d",
				translation,
				book,
				chapter,
				verse,
			)
		}
		return nil, err
	}

	return &v, nil
}

// GetPreviousVerse retrieves the previous verse in sequence
func GetPreviousVerse(translation, book string, chapter, verse int) (*Verse, error) {
	// First try to get the previous verse in the same chapter
	query := `
		SELECT id, translation, book, chapter, verse, text
		FROM bible
		WHERE translation = ? AND book = ? AND chapter = ? AND verse < ?
		ORDER BY verse DESC
		LIMIT 1
	`
	row := DB.QueryRow(query, translation, book, chapter, verse)

	var v Verse
	err := row.Scan(&v.ID, &v.Translation, &v.Book, &v.Chapter, &v.Verse, &v.Text)
	if err == nil {
		return &v, nil
	}

	if err != sql.ErrNoRows {
		return nil, err
	}

	// If no previous verse in the same chapter, try to get the last verse of the previous chapter
	query = `
		SELECT id, translation, book, chapter, verse, text
		FROM bible
		WHERE translation = ? AND book = ? AND chapter < ?
		ORDER BY chapter DESC, verse DESC
		LIMIT 1
	`
	row = DB.QueryRow(query, translation, book, chapter)
	err = row.Scan(&v.ID, &v.Translation, &v.Book, &v.Chapter, &v.Verse, &v.Text)
	if err == nil {
		return &v, nil
	}

	if err != sql.ErrNoRows {
		return nil, err
	}

	// If no previous chapter in the same book, try to get the last verse of the previous book
	query = `
		SELECT id, translation, book, chapter, verse, text
		FROM bible
		WHERE translation = ? AND book < ?
		ORDER BY book DESC, chapter DESC, verse DESC
		LIMIT 1
	`
	row = DB.QueryRow(query, translation, book)
	err = row.Scan(&v.ID, &v.Translation, &v.Book, &v.Chapter, &v.Verse, &v.Text)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf(
				"no previous verse found: %s %s %d:%d",
				translation,
				book,
				chapter,
				verse,
			)
		}
		return nil, err
	}

	return &v, nil
}

// SearchVerses searches for verses containing the given text
func SearchVerses(translation, searchText string) ([]*Verse, error) {
	query := `
		SELECT id, translation, book, chapter, verse, text
		FROM bible
		WHERE translation = ? AND (book LIKE ? OR text LIKE ?)
		ORDER BY book, chapter, verse
		LIMIT 50
	`
	searchPattern := "%" + searchText + "%"
	rows, err := DB.Query(query, translation, searchPattern, searchPattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var verses []*Verse
	for rows.Next() {
		var v Verse
		if err := rows.Scan(&v.ID, &v.Translation, &v.Book, &v.Chapter, &v.Verse, &v.Text); err != nil {
			return nil, err
		}
		verses = append(verses, &v)
	}

	return verses, nil
}

// ParseBibleReference parses a Bible reference string (e.g., "John 3:16")
// and returns the book, chapter, and verse
func ParseBibleReference(reference string) (string, int, int, error) {
	parts := strings.Split(reference, " ")
	if len(parts) < 2 {
		return "", 0, 0, fmt.Errorf("invalid Bible reference format: %s", reference)
	}

	book := strings.Join(parts[:len(parts)-1], " ")

	if strings.Contains(book, "1") && !strings.Contains(book, "1st") {
		book = strings.Replace(book, "1", "1st", 1)
	} else if strings.Contains(book, "2") && !strings.Contains(book, "2nd") {
		book = strings.Replace(book, "2", "2nd", 1)
	} else if strings.Contains(book, "3") && !strings.Contains(book, "3rd") {
		book = strings.Replace(book, "3", "3rd", 1)
	}

	chapterVerse := parts[len(parts)-1]

	cvParts := strings.Split(chapterVerse, ":")
	if len(cvParts) != 2 {
		return "", 0, 0, fmt.Errorf("invalid chapter:verse format: %s", chapterVerse)
	}

	chapter, err := parseIntWithError(cvParts[0], "chapter")
	if err != nil {
		return "", 0, 0, err
	}

	verse, err := parseIntWithError(cvParts[1], "verse")
	if err != nil {
		return "", 0, 0, err
	}

	return book, chapter, verse, nil
}

// parseIntWithError parses a string to an integer with a descriptive error
func parseIntWithError(s, field string) (int, error) {
	var result int
	_, err := fmt.Sscanf(s, "%d", &result)
	if err != nil {
		return 0, fmt.Errorf("invalid %s number: %s", field, s)
	}
	return result, nil
}

// SeedBibleData loads Bible data from JSON files in the data directory
func SeedBibleData() error {
	// Get all JSON files in the data directory
	files, err := filepath.Glob("./data/*.json")
	if err != nil {
		return err
	}

	for _, file := range files {
		// Extract translation name from filename
		translation := strings.TrimSuffix(filepath.Base(file), ".json")

		// Check if this translation is already in the database
		var count int
		err := DB.QueryRow("SELECT COUNT(*) FROM bible WHERE translation = ? LIMIT 1", translation).
			Scan(&count)
		if err != nil {
			return err
		}

		// Skip if already seeded
		if count > 0 {
			log.Printf("Translation %s already seeded, skipping...\n", translation)
			continue
		}

		// Load and parse the JSON file
		log.Printf("Seeding %s translation from %s...\n", translation, file)
		err = loadBibleFile(file, translation)
		if err != nil {
			return fmt.Errorf("error loading %s: %w", file, err)
		}
	}

	return nil
}

// loadBibleFile loads a single Bible JSON file into the database
func loadBibleFile(filePath, translation string) error {
	// Read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Parse the JSON
	var bibleData BibleData
	err = json.Unmarshal(data, &bibleData)
	if err != nil {
		return err
	}

	// Begin a transaction for faster inserts
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Prepare the insert statement
	stmt, err := tx.Prepare(`
		INSERT INTO bible (translation, book, chapter, verse, text)
		VALUES (?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Insert each verse
	for book, chapters := range bibleData.Books {
		for chapterStr, chapterData := range chapters {
			// Parse chapter number
			chapter, err := parseIntWithError(chapterStr, "chapter")
			if err != nil {
				continue // Skip invalid entries
			}

			// Insert each verse in the chapter
			for verseStr, text := range chapterData.Verses {
				// Parse verse number
				verse, err := parseIntWithError(verseStr, "verse")
				if err != nil {
					continue // Skip invalid entries
				}

				// Insert the verse
				_, err = stmt.Exec(translation, book, chapter, verse, text)
				if err != nil {
					return err
				}
			}
		}
	}

	// Commit the transaction
	return tx.Commit()
}

// SeedChapterHeaders reads JSON files and stores per-chapter headers.
// Safe to run multiple times thanks to INSERT OR IGNORE.
func SeedChapterHeaders() error {
	files, err := filepath.Glob("./data/*.json")
	if err != nil {
		return err
	}

	for _, file := range files {
		translation := strings.TrimSuffix(filepath.Base(file), ".json")

		data, readErr := os.ReadFile(file)
		if readErr != nil {
			return readErr
		}

		var bibleData BibleData
		if unmarshalErr := json.Unmarshal(data, &bibleData); unmarshalErr != nil {
			return unmarshalErr
		}

		tx, beginErr := DB.Begin()
		if beginErr != nil {
			return beginErr
		}
		var txErr error
		defer func() {
			if txErr != nil {
				tx.Rollback()
			}
		}()

		stmt, prepErr := tx.Prepare(`
			INSERT OR IGNORE INTO chapter_headers (translation, book, chapter, header)
			VALUES (?, ?, ?, ?)
		`)
		if prepErr != nil {
			tx.Rollback()
			return prepErr
		}
		defer stmt.Close()

		for book, chapters := range bibleData.Books {
			for chapterStr, chapterData := range chapters {
				if chapterData.Header == "" {
					continue
				}
				chapterNum, parseErr := parseIntWithError(chapterStr, "chapter")
				if parseErr != nil {
					continue
				}
				if _, execErr := stmt.Exec(translation, book, chapterNum, chapterData.Header); execErr != nil {
					txErr = execErr
					tx.Rollback()
					return execErr
				}
			}
		}

		if commitErr := tx.Commit(); commitErr != nil {
			return commitErr
		}
		log.Printf("Seeded headers for %s from %s...\n", translation, file)
	}

	return nil
}
