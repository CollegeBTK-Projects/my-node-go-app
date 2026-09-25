package handlers

import (
	"net/http"

	"github.com/CollegeBTK-Projects/my-node-go-app/middleware"
)

func TestError(w http.ResponseWriter, r *http.Request) {
	panic(&middleware.AppError{
		Status:  http.StatusInternalServerError,
		Message: "Тестовая ошибка с маршрута /error",
	})
}

func TestAsyncError(w http.ResponseWriter, r *http.Request) {
	errCh := make(chan error, 1)

	go func() {
		errCh <- &middleware.AppError{
			Status:  http.StatusInternalServerError,
			Message: "Асинхронная ошибка из горутины (/async-error)",
		}
	}()

	if err := <-errCh; err != nil {
		panic(err) 
	}
}
