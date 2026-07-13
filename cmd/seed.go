package main

import (
	"database/sql"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"

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
	CREATE TABLE IF NOT EXISTS movie (
		id 				TEXT PRIMARY KEY,
		title 			TEXT NOT NULL,
		poster 			TEXT NOT NULL,
		rows 			INTEGER NOT NULL,
		seats_per_row 	INTEGER NOT NULL
	);
	CREATE TABLE IF NOT EXISTS admin (
		username 	TEXT PRIMARY KEY,
		password 	TEXT NOT NULL
	);
	`); err != nil {
		return err
	}

	return nil
}

func SeedMovieDB(db *sql.DB) error {
	godotenv.Load()
	adminUsername := os.Getenv("ADMIN_USERNAME")
	adminPassword := os.Getenv("ADMIN_PASSWORD")

	if adminUsername == "" || adminPassword == "" {
		log.Print("either ADMIN_USERNAME or ADMIN_PASSWORD .env is empty")
		log.Print("[WARNING]: resorting to default ADMIN_USERNAME=admin123 and ADMIN_PASSWORD=admin123")
		adminUsername = "admin123"
		adminPassword, _ = HashPassword("admin123") // i know the password is obvious but it's just best practice okay??
	} else {
		adminPassword, _ = HashPassword(adminPassword)
	}

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

	stmt2, err := db.Prepare(`INSERT INTO admin (username, password) VALUES (?, ?)`)
	if err != nil {
		return err
	}
	defer stmt2.Close()

	stmt2.Exec(adminUsername, adminPassword)

	return tx.Commit()
}

// TODO: NOTE: we should probably moves this to utils
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
