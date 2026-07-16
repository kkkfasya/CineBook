package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
	"time"
)

func WriteJson(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJson(w, status, map[string]string{"error": err.Error()})
}

func RedisPing(client *redis.Client, maxAttempts int, backoffSec time.Duration) error {
	if maxAttempts < 1 {
		maxAttempts = 1
	}
	var lastErr error

	for attempt := range maxAttempts {
		ctx, cancel := context.WithTimeout(context.Background(), backoffSec*time.Second)
		err := client.Ping(ctx).Err()
		cancel()

		if err == nil {
			if attempt > 0 {
				log.Printf("Redis ping succeeded after %d attempts", attempt+1)
			}
			return nil
		}

		lastErr = err
		log.Printf("Redis ping attempt %d/%d failed: %v", attempt+1, maxAttempts, err)

		if attempt < maxAttempts-1 {
			backoff := time.Duration(400*1<<uint(attempt)) * time.Millisecond
			time.Sleep(backoff)
		}
	}

	return fmt.Errorf("redis ping failed after %d attempts: %w", maxAttempts, lastErr)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GetAdminCredsEnv() (adminUsername, adminPassword string) {
	godotenv.Load()
	adminUsername = os.Getenv("ADMIN_USERNAME")
	adminPassword = os.Getenv("ADMIN_PASSWORD")

	if adminUsername == "" || adminPassword == "" {
		log.Print("either ADMIN_USERNAME or ADMIN_PASSWORD .env is empty")
		log.Print("[WARNING]: resorting to default ADMIN_USERNAME=admin123 and ADMIN_PASSWORD=admin123")
		return "admin123", "admin123"
	}
	return adminUsername, adminPassword
}
