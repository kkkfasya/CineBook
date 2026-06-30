package booking

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/kkkfasya/CineBook/internal/utils"
)

// TODO:show film poster to FE

type handler struct {
	svc *Service
}

type seatInfo struct {
	SeatID    string `json:"seat_id"`
	UserID    string `json:"user_id"`
	Booked    bool   `json:"booked"`
}

func NewHandler(svc *Service) *handler {
	return &handler{svc: svc}
}

// TODO: handle err
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
		log.Println(err)
		return
	}

	data := Booking{
		UserID:  req.UserID,
		SeatID:  seatID,
		MovieID: movieID,
	}

	session, err := h.svc.Book(data)

	if err != nil {
		log.Println(err)
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

	// the FE only cares about this data
	for _, b := range bookings {
		seats = append(seats, seatInfo{
			SeatID:    b.SeatID,
			UserID:    b.UserID,
			Booked:    true,
		})

	}

	utils.WriteJson(w, http.StatusOK, seats)
}
