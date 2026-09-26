package middleware

import (
	"encoding/json"
	"log"
	"net/http"
)

type AppError struct {
	Status  int
	Message string
}

func (e *AppError) Error() string { return e.Message }

func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				if appErr, ok := err.(*AppError); ok {
					w.Header().Set("Content-Type", "application/json; charset=utf-8")
					w.WriteHeader(appErr.Status)
					json.NewEncoder(w).Encode(map[string]any{
						"error":  appErr.Message,
						"status": appErr.Status,
					})
					return
				}

				log.Printf("Паника: %v", err)
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]any{
					"error":  "Внутренняя ошибка сервера",
					"status": 500,
				})
			}
		}()

		next.ServeHTTP(w, r)
	})
}
