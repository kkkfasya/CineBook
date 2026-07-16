package booking

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"database/sql"

	"github.com/kkkfasya/CineBook/internal/utils"
)

// TODO:show film poster to FE

// TIL this is called DTO (Data Transfer Object)
type MovieResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Poster      string `json:"poster"`
	Rows        uint8  `json:"rows"`
	SeatsPerRow uint8  `json:"seats_per_row"`
}

type handler struct {
	svc *Service
}

type seatInfo struct {
	SeatID    string `json:"seat_id"`
	UserID    string `json:"user_id"`
	Booked    bool   `json:"booked"`
	Confirmed bool   `json:"confirmed"`
}

type holdSeatRequest struct {
	UserID string `json:"user_id"`
}

type sessionResponse struct {
	SessionID string `json:"session_id"`
	MovieID   string `json:"movie_id"`
	SeatID    string `json:"seat_id"`
	UserID    string `json:"user_id"`
	Status    string `json:"status"`
	ExpiresAt string `json:"expires_at,omitempty"`
}

func NewHandler(svc *Service) *handler {
	return &handler{svc: svc}
}

func (h *handler) HoldSeat(w http.ResponseWriter, r *http.Request) {
	type holdPayloadRequest struct {
		UserID string `json:"user_id"`
	}

	type holdResponse struct {
		SessionID string `json:"session_id"`
		MovieID   string `json:"movie_id"`
		SeatID    string `json:"seat_id"`
		ExpiresAt string `json:"expires_at"`
	}

	movieID := r.PathValue("movieID")
	seatID := r.PathValue("seatID")

	var req holdPayloadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	data := Booking{
		UserID:  req.UserID,
		SeatID:  seatID,
		MovieID: movieID,
	}

	session, err := h.svc.Book(data)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, holdResponse{
		SeatID:    seatID,
		SessionID: session.ID,
		MovieID:   session.MovieID,
		ExpiresAt: session.ExpiresAt.Format(time.RFC3339),
	})
}

func (h *handler) ListSeats(w http.ResponseWriter, r *http.Request) {
	movieID := r.PathValue("movieID")
	bookings := h.svc.ListBookings(movieID)
	seats := make([]seatInfo, 0, len(bookings))

	for _, b := range bookings {
		seats = append(seats, seatInfo{
			SeatID:    b.SeatID,
			UserID:    b.UserID,
			Booked:    true,
			Confirmed: b.Status == "confirmed",
		})

	}

	utils.WriteJson(w, http.StatusOK, seats)
}

func (h *handler) ConfirmSession(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("sessionID")

	var req holdSeatRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if req.UserID == "" {
		utils.WriteError(w, http.StatusBadRequest, ErrMissingUserID)
		return
	}

	session, err := h.svc.ConfirmSeat(r.Context(), sessionID, req.UserID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, sessionResponse{
		SessionID: session.ID,
		MovieID:   session.MovieID,
		SeatID:    session.SeatID,
		UserID:    req.UserID,
		Status:    session.Status,
	})

}

func (h *handler) ReleaseSession(w http.ResponseWriter, r *http.Request) {
	var req holdSeatRequest
	sid := r.PathValue("sessionID")

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if req.UserID == "" {
		utils.WriteError(w, http.StatusBadRequest, ErrMissingUserID)
		return
	}

	err := h.svc.ReleaseSeat(r.Context(), sid, req.UserID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) ListMovies(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		movies := []MovieResponse{}
		rows, err := db.QueryContext(r.Context(), "SELECT * FROM movie") // always use QueryContext in http request setting
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}
		defer rows.Close()

		// cursor pointing just before the first row of data. thus we need to call Next()  to move the cursor and check if data is available
		for rows.Next() {
			var m MovieResponse
			if err := rows.Scan(&m.ID, &m.Title, &m.Poster, &m.Rows, &m.SeatsPerRow); err != nil {
				utils.WriteError(w, http.StatusInternalServerError, err)
				return
			}
			movies = append(movies, m)
		}

		if err := rows.Err(); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}
		log.Print(movies)
		utils.WriteJson(w, http.StatusOK, movies)
	})
}
