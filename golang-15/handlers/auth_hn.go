package handlers

import "net/http"

func Auth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Вы админ!"))
}
