package middleware

import (
	"log"
	"net/http"
	"time"
)

func Logging(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		duration := time.Since(start).Milliseconds()

		log.Printf("[%s] %s %s - %dms\n", start.Format("2006-01-02 15:04:05"), r.Method, r.URL.Path, duration)
	})
}
