package booking

type MemoryStore struct {
	bookings map[string]Booking // go maps are't concurrent safe, race-condition is very possible
}

// this is like constructor in OOP language
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		bookings: map[string]Booking{}, // dont forget to initialize
	}
}

func (m *MemoryStore) Book(b Booking) error {
	if _, exists := m.bookings[b.SeatID]; exists {
		return ErrSeatAlreadyBooked
	}

	// use ID perhaps?
	m.bookings[b.SeatID] = b
	return nil
}

func (m *MemoryStore) ListBookings(movieID string) []Booking {
	var r []Booking

	if len(m.bookings) == 0 {
		return r
	}

	// O(n), TODO: try to implement O(1) version
	for _, b := range m.bookings {
		if b.MovieID == movieID {
			r = append(r, b)
		}
	}

	return r
}
