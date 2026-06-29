package booking

import (
	"net/http"

	"github.com/kkkfasya/CineBook/internal/utils"
)

// TODO:show film poster to FE

type handler struct {
	svc *Service
}

type seatInfo struct {
	SeatID string `json:"seat_id"`
	UserID string `json:"user_id"`
	Booked bool   `json:"booked"`
}

func NewHandler(svc *Service) *handler {
	return &handler{svc: svc}
}

func (h *handler) ListSeats(w http.ResponseWriter, r *http.Request) {
	movieID := r.PathValue("movieID")
	bookings := h.svc.ListBookings(movieID)
	seats := make([]seatInfo, 0, len(bookings))

	// the FE only cares about this data
	for _, b := range bookings {
		seats = append(seats, seatInfo{
			SeatID: b.SeatID,
			UserID: b.UserID,
			Booked: true,
		})

	}

	utils.WriteJson(w, http.StatusOK, seats)
}
