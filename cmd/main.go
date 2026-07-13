package main

import (
	"log"
	"net/http"

	"github.com/redis/go-redis/v9"

	"github.com/kkkfasya/CineBook/internal/booking"
	mw "github.com/kkkfasya/CineBook/internal/middleware"
	"github.com/kkkfasya/CineBook/internal/utils"

	"database/sql"
	_ "github.com/ncruces/go-sqlite3/driver"
)

const PORT = ":8080"

func main() {
	db, err := sql.Open("sqlite3", "movies.db")
	CreateMovieDB(db)

	if err != nil {
		log.Fatal(err)
	}

	rclient := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	if err := utils.RedisPing(rclient, 3, 3); err != nil {
		log.Fatal(err)
	}
	store := booking.NewRedisStore(rclient)
	svc := booking.NewService(store)
	handler := booking.NewHandler(svc)

	mux := http.NewServeMux()

	//TODO: use go embed to serve this static file
	// https://oneuptime.com/blog/post/2026-01-25-bundle-static-assets-go-embed/view
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("GET /", fs)

	// TODO: add redirect for normal user here
	mux.HandleFunc("GET /admin", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/admin.html")
	})
	mux.HandleFunc("GET /login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/login.html")
	})

	mux.HandleFunc("GET /api/v1/movies", listMovies)
	mux.HandleFunc("GET /api/v1/movies/{movieID}/seats", handler.ListSeats)
	mux.HandleFunc("POST /api/v1/movies/{movieID}/seats/{seatID}/hold", handler.HoldSeat)
	mux.HandleFunc("PUT /api/v1/sessions/{sessionID}/confirm", handler.ConfirmSession)
	mux.HandleFunc("DELETE /api/v1/sessions/{sessionID}", handler.ReleaseSession)

	log.Printf("server started at http://localhost%s\n", PORT)
	if err := http.ListenAndServe(PORT, mw.StripTrailingSlash(mux)); err != nil {
		log.Printf("server at localhost:%s failed\n", PORT)
		log.Fatal(err)
	}

}

// yes we will hardcode it as for now
var movies = []MovieResponse{
	{ID: "cb", Title: "Call Boy", Rows: 3, SeatsPerRow: 3},
	{ID: "mhs", Title: "Un homme qui dort", Rows: 6, SeatsPerRow: 6},
}

func listMovies(w http.ResponseWriter, r *http.Request) {
	utils.WriteJson(w, http.StatusOK, movies)
}
