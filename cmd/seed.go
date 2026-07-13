package main

import (
	"crypto/rand"
	"database/sql"
	"time"

	_ "github.com/ncruces/go-sqlite3/driver"

	"github.com/joho/godotenv"
	"github.com/oklog/ulid/v2"
)

type MovieResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Poster      string `json:"poster"`
	Rows        uint8  `json:"rows"`
	SeatsPerRow uint8  `json:"seats_per_row"`
}

func CreateMovieDB(db *sql.DB) error {
	if _, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS movies (
		id 				TEXT PRIMARY KEY,
		title 			TEXT NOT NULL,
		poster 			TEXT NOT NULL,
		rows 			INTEGER NOT NULL,
		seats_per_row 	INTEGER NOT NULL
	);
	CREATE TABLE IF NOT EXISTS admins (
		username TEXT PRIMARY KEY,
		password_hash TEXT NOT NULL
	);
	`); err != nil {
		return err
	}

	return nil
}

func SeedMovieDB(db *sql.DB) error {
	entropy := rand.Reader
	id, err := ulid.New(ulid.Timestamp(time.Now()), entropy)
	if err != nil {
		return err
	}
	// TODO: prepare 3 films, id is ulid, poster get from TMDb
	if _, err := db.Exec(
		`
	INSERT INTO movies 
	(id, title, poster, rows, seats_per_row)
	VALUES
			('Blue Train', 'John Coltrane', 56.99),
			('Giant Steps', 'John Coltrane', 63.99),
			('Jeru', 'Gerry Mulligan', 17.99),
			('Sarah Vaughan', 'Sarah Vaughan', 34.98)
	`); err != nil {
		return err
	}
	return nil
}
