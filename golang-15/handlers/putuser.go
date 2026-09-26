package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/CollegeBTK-Projects/my-node-go-app/models"
)

func PutUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if id == "" {
		http.Error(w, "ID не указан", http.StatusBadRequest)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "ID должен быть числом", http.StatusBadRequest)
		return
	}

	list := []models.User{}
	u := models.User{
		ID: idInt,
	}

	err = json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	readfile, err := os.Open("users.json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewDecoder(readfile).Decode(&list)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	find := false
	for i := range list {
		if list[i].ID == idInt {
			list[i] = u
			find = true
		}
	}
	if !find {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	file, err := os.Create("users.json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(list)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
