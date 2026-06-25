package booking

import (
	"errors"
)

var (
	ErrSeatAlreadyBooked = errors.New("seat already taken")
)

type Booking struct {
	ID      string
	MovieID string
	SeatID  string
	UserID  string
	Status  string // can we use enum for this?
}

type BookingStore interface {
	Book(b Booking) error
	ListBookings(movieID string) []Booking
}
