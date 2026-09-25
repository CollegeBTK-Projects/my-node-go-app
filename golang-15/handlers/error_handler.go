package handlers

import (
	"errors"
	"net/http"
)

func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	panic(errors.New("Ошибка"))
}
