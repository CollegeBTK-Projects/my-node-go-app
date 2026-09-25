package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/CollegeBTK-Projects/my-node-go-app/models"
)

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	file, err := os.OpenFile("users.json", os.O_CREATE|os.O_RDONLY, 0o664)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	var u []models.User

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	js := json.NewDecoder(file)
	if err := js.Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(u) == 0 {
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("Пусто..."))
		return
	}

	w.WriteHeader(http.StatusOK)
	for _, v := range u {
		fmt.Fprintf(w, "ID: %d User: %s Group: %s\n", v.ID, v.Name, v.Group)
	}
}
