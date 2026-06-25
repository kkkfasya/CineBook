package booking

type Booking struct {
	ID      string
	MovieID string
	SeatID  string
	UserID  string
	Status  string // can we use enum for this?
}

type BookingStore interface {
	Book(b Booking) (Booking, error)
	ListBookings(movieID string) []Booking
}
