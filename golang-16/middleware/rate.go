package middleware

import (
	"encoding/json"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type RateLimiter struct {
	mu       sync.Mutex
	requests map[string]int
	limit    int
	resetAt  time.Time
}

func NewRateLimiter(limit int) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string]int),
		limit:    limit,
		resetAt:  time.Now().Add(time.Minute),
	}
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)

		rl.mu.Lock()
		if time.Now().After(rl.resetAt) {
			rl.requests = make(map[string]int)
			rl.resetAt = time.Now().Add(time.Minute)
		}
		rl.requests[ip]++
		count := rl.requests[ip]
		rl.mu.Unlock()

		w.Header().Set("X-RateLimit-Limit", strconv.Itoa(rl.limit))
		w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(rl.limit-count))

		if count > rl.limit {
			w.Header().Set("Retry-After", "60")
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]any{
				"error":  "Слишком много запросов",
				"status": 429,
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}
