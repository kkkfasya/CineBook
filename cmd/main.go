package main

import (
	"log"
	"net/http"
	"os"

	"github.com/redis/go-redis/v9"

	"github.com/kkkfasya/CineBook/internal/booking"
	mw "github.com/kkkfasya/CineBook/internal/middleware"
	"github.com/kkkfasya/CineBook/internal/utils"

	"database/sql"

	_ "github.com/ncruces/go-sqlite3/driver"
)

const PORT = ":8080"
const DB_NAME = "movies.db"

func main() {
	os.Remove("movies.db") // in-memory does not work somehow
	db, err := sql.Open("sqlite3", "movies.db")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	if err := CreateMovieDB(db); err != nil {
		log.Fatal(err)
	}
	if err := SeedMovieDB(db); err != nil {
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
	// authGuard := mw.RequiresAdmin()
	basicAuth := mw.BasicAuth(utils.GetAdminCredsEnv())

	//TODO: use go embed to serve this static file
	// https://oneuptime.com/blog/post/2026-01-25-bundle-static-assets-go-embed/view
	mux.Handle("GET /", http.FileServer(http.Dir("static")))

	mux.Handle("GET /admin", basicAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/admin.html")
	})))

	mux.HandleFunc("GET /login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/login.html")
	})

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mux.Handle("GET /api/v1/movies", handler.ListMovies(db))

	mux.HandleFunc("GET /api/v1/movies/{movieID}/seats", handler.ListSeats)
	mux.HandleFunc("POST /api/v1/movies/{movieID}/seats/{seatID}/hold", handler.HoldSeat)
	mux.HandleFunc("PUT /api/v1/sessions/{sessionID}/confirm", handler.ConfirmSession)
	mux.HandleFunc("DELETE /api/v1/sessions/{sessionID}", handler.ReleaseSession)

	// admin movies CRUD
	mux.Handle("GET /api/v1/admin/movies", handler.ListMovies(db))
	mux.HandleFunc("POST /api/v1/admin/movies", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	})
	mux.HandleFunc("PUT /api/v1/admin/movies/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	})
	mux.HandleFunc("DELETE /api/v1/admin/movies/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	})

	log.Printf("server started at http://localhost%s\n", PORT)
	if err := http.ListenAndServe(PORT, mw.StripTrailingSlash(mux)); err != nil {
		log.Printf("server at localhost:%s failed\n", PORT)
		log.Fatal(err)
	}
}
