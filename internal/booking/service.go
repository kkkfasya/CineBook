package booking

type Service struct {
	store BookingStore
}

func NewService(store BookingStore) *Service {
	return &Service{store}
}

func (s *Service) Book(b Booking) (Booking, error) {
	return s.store.Book(b)
}

// TODO: understand this fn when im not sleepy
func (s *Service) ListBookings(movieID string) []Booking {
	return s.store.ListBookings(movieID)
}
