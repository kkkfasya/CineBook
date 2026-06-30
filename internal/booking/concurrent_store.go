package booking

import "sync"

// use Mutex for locking mechanism to prevent race condition
// better yet use RWMutex if it's read-heavy or need multiple reader but one writer
type ConcurrentStore struct {
	bookings map[string]Booking
	sync.RWMutex
}

func NewConcurrentStore() *ConcurrentStore {
	return &ConcurrentStore{
		bookings: map[string]Booking{},
	}
}

func (m *ConcurrentStore) Book(b Booking) error {
	m.Lock()
	defer m.Unlock()

	if _, exists := m.bookings[b.SeatID]; exists {
		return ErrSeatAlreadyBooked
	}

	m.bookings[b.SeatID] = b
	return nil
}

func (m *ConcurrentStore) ListBookings(movieID string) []Booking {
	m.RLock()
	defer m.RUnlock()
	var r []Booking

	if len(m.bookings) == 0 {
		return r
	}

	for _, b := range m.bookings {
		if b.MovieID == movieID {
			r = append(r, b)
		}
	}

	return r
}
