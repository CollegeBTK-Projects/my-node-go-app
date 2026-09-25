package middleware

import (
    "log"
    "net/http"
    "time"
)

type statusRecorder struct {
    http.ResponseWriter
    status int
}

func (r *statusRecorder) WriteHeader(status int) {
    r.status = status
    r.ResponseWriter.WriteHeader(status)
}

func Logger(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
        next.ServeHTTP(rec, r)

        log.Printf("[%s] %s %s %d - %s",
            start.Format("2006-01-02 15:04:05"),
            r.Method,
            r.URL.Path,
            rec.status,
            time.Since(start).Round(time.Millisecond),
        )
    })
}