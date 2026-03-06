package ratelimit

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
)

var (
	redisClient *redis.Client
	limiter     *redis_rate.Limiter
)

func InitRateLimiter(redisAddr string) error {
	// Connexion Redis
	redisClient = redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("Unable to connect to Redis:: %w", err)
	}

	log.Println("Connection to Redis successful")

	limiter = redis_rate.NewLimiter(redisClient)

	return nil
}

func RateLimitMiddleware(requestsPerMinute int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := getClientIP(r)

			key := fmt.Sprintf("rate_limit:%s", ip)

			res, err := limiter.Allow(r.Context(), key, redis_rate.PerMinute(requestsPerMinute))
			if err != nil {
				log.Printf("Erreur rate limit: %v", err)
				http.Error(w, "Error Server", http.StatusInternalServerError)
				return
			}

			if res.Allowed == 0 {
				w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", requestsPerMinute))
				w.Header().Set("X-RateLimit-Remaining", "0")
				w.Header().Set("Retry-After", fmt.Sprintf("%d", res.RetryAfter/time.Second))

				http.Error(w, "Too many attempts. Please try again later.", http.StatusTooManyRequests)
				return
			}

			w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", requestsPerMinute))
			w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", res.Remaining))

			next.ServeHTTP(w, r)
		})
	}
}

func getClientIP(r *http.Request) string {
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		ips := strings.Split(forwarded, ",")
		return strings.TrimSpace(ips[0])
	}

	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	ip := r.RemoteAddr
	if colon := strings.LastIndex(ip, ":"); colon != -1 {
		ip = ip[:colon]
	}

	return ip
}

func CloseRedis() error {
	if redisClient != nil {
		return redisClient.Close()
	}
	return nil
}
