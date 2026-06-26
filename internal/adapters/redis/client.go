package redis // NOTE: should we really name it this?

import (
	"context"
	"log"

	goredis "github.com/redis/go-redis/v9"
)

func NewRedisClient(addr string) *goredis.Client {
	rdb := goredis.NewClient(&goredis.Options{Addr: addr})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("redis ping: %v", err)
	}
	log.Printf("connected to redis at %s", addr)

	return rdb
}
