package booking

import (
	"sync"
)

// use Mutex for locking mechanism to prevent race condition
// better yet use RWMutex if it's read-heavy or need multiple reader but one writer
type ConcurrentStore struct {
	sync.RWMutex
	bookings map[string]Booking
	byMovie  map[string]map[string]struct{} // we use nested struct because removing it is O(1) instead of O(n)
}

func NewConcurrentStore() *ConcurrentStore {
	return &ConcurrentStore{
		bookings: map[string]Booking{},
		byMovie:  map[string]map[string]struct{}{},
	}
}

func (m *ConcurrentStore) Book(b Booking) error {
	m.Lock()
	defer m.Unlock()

	if _, exists := m.bookings[b.SeatID]; exists {
		return ErrSeatAlreadyBooked
	}

	m.bookings[b.SeatID] = b

	if m.byMovie[b.MovieID] == nil {
		m.byMovie[b.MovieID] = map[string]struct{}{}
	}

	m.byMovie[b.MovieID][b.SeatID] = struct{}{}

	return nil
}

func (m *ConcurrentStore) ListBookings(movieID string) []Booking {
	m.RLock()
	defer m.RUnlock()

	seatIDs, exists := m.byMovie[movieID]
	if !exists || len(seatIDs) == 0 {
		return nil
	}

	r := make([]Booking, 0, len(seatIDs))

	for sid := range seatIDs {
		if b, found := m.bookings[sid]; found {
			r = append(r, b)
		}
	}

	return r

}
