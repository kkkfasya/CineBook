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

const holdTTL = 1 * time.Minute

type RedisStore struct {
	rdb *redis.Client
}

func NewRedisStore(rdb *redis.Client) *RedisStore {
	return &RedisStore{rdb: rdb}
}

func sessionKey(id string) string {
	return fmt.Sprintf("session:%s", id)
}

func seatKey(movieID, seatID string) string {
	return fmt.Sprintf("seat:%s:%s", movieID, seatID)
}

func parseSession(s string) (Booking, error) {

	var data Booking
	if err := json.Unmarshal([]byte(s), &data); err != nil {
		return Booking{}, err
	}

	return Booking{
		ID:      data.ID,
		MovieID: data.MovieID,
		SeatID:  data.SeatID,
		UserID:  data.UserID,
		Status:  data.Status,
	}, nil
}

// key & value
//
// sessionkey -> seatkey
//
// seatkey -> booking stuct
func (r *RedisStore) hold(b Booking) (Booking, error) {

	id := uuid.New().String()
	b.ID = id

	now := time.Now()
	ctx := context.Background()
	seatKey := seatKey(b.MovieID, b.SeatID)
	bookingVal, err := json.Marshal(b)

	if err != nil {
		return Booking{}, err
	}

	if cmd := r.rdb.SetArgs(ctx, sessionKey(id), seatKey, redis.SetArgs{TTL: holdTTL}); cmd.Val() != "OK" {
		log.Print(cmd.Val())
		return Booking{}, ErrFailedToSetSessionKey
	}

	if cmd := r.rdb.SetArgs(ctx, seatKey, bookingVal, redis.SetArgs{
		Mode: "NX", // NX set if key not exist, otherwise ignore
		TTL:  holdTTL,
	}); cmd.Val() != "OK" {
		log.Print(cmd.Val())
		return Booking{}, ErrSeatAlreadyBooked
	}

	return Booking{
		ID:        id,
		MovieID:   b.MovieID,
		SeatID:    b.SeatID,
		UserID:    b.UserID,
		Status:    "held",
		ExpiresAt: now.Add(holdTTL),
	}, nil
}

func (r *RedisStore) Book(b Booking) (Booking, error) {
	session, err := r.hold(b)
	if err != nil {
		return Booking{}, err
	}

	log.Printf("%s session booked: %s", session.SeatID, session)

	return session, nil
}

// TODO: add validation for userID
func (r *RedisStore) getSession(ctx context.Context, sessionID string, userID string) (Booking, string, error) {
	sk, err := r.rdb.Get(ctx, sessionKey(sessionID)).Result()

	if err != nil {
		return Booking{}, "", err
	}

	val, err := r.rdb.Get(ctx, sk).Result()
	if err != nil {
		return Booking{}, "", err
	}

	s, err := parseSession(val)
	if err != nil {
		return Booking{}, "", err
	}

	return s, sk, nil
}

func (r *RedisStore) Confirm(ctx context.Context, sessionID string, userID string) (Booking, error) {
	session, sk, err := r.getSession(ctx, sessionID, userID)
	if err != nil {
		return Booking{}, err
	}

	// persist removes TTL
	type sessionResponse struct {
		SessionID string `json:"session_id"`
		MovieID   string `json:"movie_id"`
		SeatID    string `json:"seat_id"`
		UserID    string `json:"user_id"`
		Status    string `json:"status"`
		ExpiresAt string `json:"expires_at,omitempty"`
	}
	r.rdb.Persist(ctx, sk)
	r.rdb.Persist(ctx, sessionKey(sessionID))

	session.Status = "confirmed"

	val, err := json.Marshal(Booking{
		ID:      session.ID,
		MovieID: session.MovieID,
		SeatID:  session.SeatID,
		UserID:  session.UserID,
		Status:  session.Status,
	})

	if err != nil {
		return Booking{}, err
	}
	r.rdb.Set(ctx, sk, val, 0)
	return session, nil

}

func (r *RedisStore) Release(ctx context.Context, sessionID string, userID string) error {
	_, sk, err := r.getSession(ctx, sessionID, userID)

	if err != nil {
		return err
	}

	r.rdb.Del(ctx, sk, sessionKey(sessionID))

	return nil
}

func (r *RedisStore) ListBookings(movieID string) []Booking {
	pattern := fmt.Sprintf("seat:%s:*", movieID)
	var sessions []Booking
	ctx := context.Background()

	iter := r.rdb.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		val, err := r.rdb.Get(ctx, iter.Val()).Result()
		if err != nil {
			continue
		}

		s, err := parseSession(val)
		sessions = append(sessions, s)
	}

	return sessions
}
