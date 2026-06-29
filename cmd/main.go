package main

import (
	"log"
	"net/http"

	"github.com/redis/go-redis/v9"

	"github.com/kkkfasya/CineBook/internal/booking"
	"github.com/kkkfasya/CineBook/internal/utils"
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
	store := booking.NewRedisStore(
		redis.NewClient(&redis.Options{Addr: "localhost:6379"}),
	)
	svc := booking.NewService(store)
	handler := booking.NewHandler(svc)

	//TODO: use go embed to serve this static file
	// https://oneuptime.com/blog/post/2026-01-25-bundle-static-assets-go-embed/view
	mux.Handle("GET /", http.FileServer(http.Dir("static")))
	mux.HandleFunc("GET /movies", listMovies)
	mux.HandleFunc("GET /movies/{movieID}/seats", handler.ListSeats)

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
	utils.WriteJson(w, http.StatusOK, movies)
}
