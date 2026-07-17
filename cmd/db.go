package main

import (
	"database/sql"
	_ "github.com/ncruces/go-sqlite3/driver"
	"github.com/oklog/ulid/v2"
)

func CreateMovieDB(db *sql.DB) error {
	if _, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS movie (
		id 				TEXT PRIMARY KEY,
		title 			TEXT NOT NULL,
		poster 			TEXT NOT NULL,
		rows 			INTEGER NOT NULL,
		seats_per_row 	INTEGER NOT NULL
	);
	`); err != nil {
		return err
	}

	return nil
}

func SeedMovieDB(db *sql.DB) error {
	seeds := []struct {
		Title  string
		Poster string
		Rows   int
		Seats  int
	}{
		{"Call Boy", "https://media.themoviedb.org/t/p/original/3fVRb3uoRDjbA9X95C88FOJ0rlZ.jpg", 3, 3},
		{"Un homme qui dort", "https://media.themoviedb.org/t/p/original/rhsccTC8rpKAjmSCuyATDIyKhKZ.jpg", 6, 6},
		{"The Green Ray", "https://media.themoviedb.org/t/p/original/1E3pliSC7lXWw6zJhMvG6ba0UNX.jpg", 2, 2},
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := db.Prepare(`
	INSERT INTO movie  (id, title, poster, rows, seats_per_row)
	VALUES (?, ?, ?, ?, ?)`)

	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, s := range seeds {
		id := ulid.Make().String()
		if _, err := stmt.Exec(id, s.Title, s.Poster, s.Rows, s.Seats); err != nil {
			return err
		}
	}
	return tx.Commit()
}
