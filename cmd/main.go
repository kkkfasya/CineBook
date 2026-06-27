package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type movieResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Rows        uint8  `json:"rows"`
	SeatsPerRow uint8  `json:"seats_per_row"` // NOTE: change this to columns after tutorial
}

// TODO: change any mention of the word `movie`  to `film` because i'm a cinephile
const PORT = ":8080"

func main() {
	mux := http.NewServeMux()

	//TODO: use go embed to serve this static file
	// https://oneuptime.com/blog/post/2026-01-25-bundle-static-assets-go-embed/view
	mux.Handle("GET /", http.FileServer(http.Dir("static")))

	mux.HandleFunc(
		"GET /movies",
		listMovies,
	)

	if err := http.ListenAndServe(PORT, mux); err != nil {
		log.Fatal(err)
	}
}

// yes we will hardcode it as for now
var movies = []movieResponse{
	{ID: "cb", Title: "Call Boy", Rows: 3, SeatsPerRow: 3},
	{ID: "as", Title: "A Separation", Rows: 6, SeatsPerRow: 6},
}

func listMovies(w http.ResponseWriter, r *http.Request) {
	WriteJson(w, http.StatusOK, movies)
}

func WriteJson(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
