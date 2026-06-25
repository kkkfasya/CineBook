package booking

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const holdTTL = 2 * time.Minute

type RedisStore struct {
	rdb *redis.Client
}

func NewRedisStore(rdb *redis.Client) *RedisStore {
	return &RedisStore{rdb: rdb}
}

func sessionKey(id string) string {
	return fmt.Sprintf("session:%s", id)
}

func (r *RedisStore) hold(b Booking) (Booking, error) {
	id := uuid.New().String()
	now := time.Now()

	return Booking{
		ID:        id,
		MovieID:   b.MovieID,
		SeatID:    b.SeatID,
		UserID:    b.UserID,
		Status:    "held",
		ExpiresAt: now.Add(holdTTL),
	}, nil
}

func (r *RedisStore) Book(b Booking) error {
	session, err := r.hold(b)
	if err != nil {
		return err
	}

	log.Printf("session booked: %s", session)

	return nil
}

func (m *RedisStore) ListBookings(movieID string) []Booking {
	m.RLock()
	defer m.Unlock()
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
