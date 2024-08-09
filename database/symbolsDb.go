package database

import (
	"database/sql"
	"log"
)

func Open() *sql.DB {
	db, err := sql.Open("sqlite3", "./symbols.db")

	if err != nil {
		log.Fatal("Sqlite Database", err)
	}

	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS symbols (symbol TEXT PRIMARY KEY)`)
	if err != nil {
		log.Fatal("Error while checking if table exists: ", err)
	}

	return db
}
