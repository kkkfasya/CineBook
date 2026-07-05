package booking

import (
	"context"
	"errors"
	"time"
)

var (
	ErrSeatAlreadyBooked     = errors.New("seat already taken")
	ErrFailedToSetSessionKey = errors.New("failed to set session key")
	ErrMissingUserID         = errors.New("missing user_id field")
)

type Booking struct {
	ID        string
	MovieID   string
	SeatID    string
	UserID    string
	Status    string // can we use enum for this?
	ExpiresAt time.Time
}

type BookingStore interface {
	Book(b Booking) (Booking, error)
	ListBookings(movieID string) []Booking
	Confirm(ctx context.Context, sessionID string, userID string) (Booking, error)
	Release(ctx context.Context, sessionID string, userID string) error
}
