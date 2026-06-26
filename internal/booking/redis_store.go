package booking

import (
	"context"
	"encoding/json"
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
	ctx := context.Background()
	key := fmt.Sprintf("seat:%s:%s", b.MovieID, b.SeatID)
	b.ID = id
	val, err := json.Marshal(b)

	if err != nil {
		return Booking{}, err
	}

	status := r.rdb.SetArgs(ctx, key, val, redis.SetArgs{
		Mode: "NX", // NX set if key not exist, otherwise ignore
		TTL:  holdTTL,
	})

	if status.Val() != "OK" {
		return Booking{}, ErrSeatAlreadyBooked
	}

	// XXX: is this needed since we have SetArgs already
	// perhaps i should learn more about this redis client lib
	r.rdb.Set(ctx, sessionKey(id), key, holdTTL)

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

// TODO: implement
func (r *RedisStore) ListBookings(movieID string) []Booking {
	return []Booking{}
}
