package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/kkkfasya/CineBook/internal/booking"
	"github.com/kkkfasya/CineBook/internal/utils"
)

type movieResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Rows        uint8  `json:"rows"`
	SeatsPerRow uint8  `json:"seats_per_row"`
}

// TODO: change any mention of the word `movie`  to `film` because i'm a cinephile
const PORT = ":8080"

func main() {
	rclient := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	if err := redisPing(rclient, 3, 3); err != nil {
		log.Fatal(err)
	}
	store := booking.NewRedisStore(rclient)
	svc := booking.NewService(store)
	handler := booking.NewHandler(svc)

	mux := http.NewServeMux()

	//TODO: use go embed to serve this static file
	// https://oneuptime.com/blog/post/2026-01-25-bundle-static-assets-go-embed/view
	mux.Handle("GET /", http.FileServer(http.Dir("static")))
	mux.HandleFunc("GET /api/v1/movies", listMovies)
	mux.HandleFunc("GET /api/v1/movies/{movieID}/seats", handler.ListSeats)
	mux.HandleFunc("POST /api/v1/movies/{movieID}/seats/{seatID}/hold", handler.HoldSeat)
	mux.HandleFunc("PUT /api/v1/sessions/{sessionID}/confirm", handler.ConfirmSession)
	mux.HandleFunc("DELETE /api/v1/sessions/{sessionID}", handler.ReleaseSession)

	log.Printf("server started at http://localhost%s\n", PORT)
	if err := http.ListenAndServe(PORT, mux); err != nil {
		log.Printf("server at localhost:%s failed\n", PORT)
		log.Fatal(err)
	}

}

func redisPing(client *redis.Client, maxAttempts int, backoffSec time.Duration) error {
	if maxAttempts < 1 {
		maxAttempts = 1
	}
	var lastErr error

	for attempt := range maxAttempts {
		ctx, cancel := context.WithTimeout(context.Background(), backoffSec*time.Second)
		err := client.Ping(ctx).Err()
		cancel()

		if err == nil {
			if attempt > 0 {
				log.Printf("Redis ping succeeded after %d attempts", attempt+1)
			}
			return nil
		}

		lastErr = err
		log.Printf("Redis ping attempt %d/%d failed: %v", attempt+1, maxAttempts, err)

		if attempt < maxAttempts-1 {
			backoff := time.Duration(400*1<<uint(attempt)) * time.Millisecond
			time.Sleep(backoff)
		}
	}

	return fmt.Errorf("redis ping failed after %d attempts: %w", maxAttempts, lastErr)
}

// yes we will hardcode it as for now
var movies = []movieResponse{
	{ID: "cb", Title: "Call Boy", Rows: 3, SeatsPerRow: 3},
	{ID: "as", Title: "A Separation", Rows: 6, SeatsPerRow: 6},
}

func listMovies(w http.ResponseWriter, r *http.Request) {
	utils.WriteJson(w, http.StatusOK, movies)
}
