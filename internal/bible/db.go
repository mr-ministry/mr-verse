package bible

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// DB is the global database connection
var DB *sql.DB

// getDBPath returns the path to the SQLite database file
// It checks for the DB_PATH environment variable and falls back to a default
func getDBPath() string {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		// Allow app to access db on the same dir as the executable
		// if there is no .env or the DB_PATH varaible is not set
		dbPath = "./bible.db"
	}
	return dbPath
}

// OpenDB opens a connection to the SQLite database
// and returns a pointer to the database
// and an error if it occurs.
func OpenDB() (*sql.DB, error) {
	// Get the database path
	dbPath := getDBPath()

	// Ensure the data directory exists
	err := os.MkdirAll(filepath.Dir(dbPath), 0755)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

// InitDB initializes the database connection and creates the tables if they don't exist
func InitDB() error {
	var err error
	DB, err = OpenDB()
	if err != nil {
		return err
	}

	// Create the tables if they don't exist
	err = createTables()
	if err != nil {
		return err
	}

	return nil
}

// createTables creates the necessary tables in the database
func createTables() error {
	// Create the bible table
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS bible (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			translation TEXT NOT NULL,
			book TEXT NOT NULL,
			chapter INTEGER NOT NULL,
			verse INTEGER NOT NULL,
			text TEXT NOT NULL,
			UNIQUE(translation, book, chapter, verse)
		)
	`)
	if err != nil {
		return err
	}

	// Create index for faster lookups
	_, err = DB.Exec(`
		CREATE INDEX IF NOT EXISTS idx_bible_lookup 
		ON bible(translation, book, chapter, verse)
	`)
	if err != nil {
		return err
	}

	return nil
}

// CloseDB closes the database connection
func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
