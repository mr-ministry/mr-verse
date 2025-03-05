package bible

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// OpenDB opens a connection to the SQLite database
// and returns a pointer to the database
// and an error if it occurs.
func OpenDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./data/bible.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	return db, nil
}
