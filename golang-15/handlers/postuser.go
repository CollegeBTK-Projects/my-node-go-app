package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/CollegeBTK-Projects/my-node-go-app/models"
)

func GetLastID() (int, error) {
	file, err := os.Open("users.json")
	if err != nil {
		return 0, err
	}
	defer file.Close()

	var u []models.User
	js := json.NewDecoder(file)
	if err := js.Decode(&u); err != nil {
		if err == io.EOF {
			return 0, nil
		}
		return 0, err
	}
	var lastID int
	if len(u) > 0 {
		lastID = u[len(u)-1].ID
	} else {
		lastID = 1
	}
	return lastID, nil
}

func PostUser(w http.ResponseWriter, r *http.Request) {
	num, err := GetLastID()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	list := []models.User{}
	u := models.User{
		ID: num + 1,
	}

	readfile, err := os.Open("users.json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer readfile.Close()

	err = json.NewDecoder(readfile).Decode(&list)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	list = append(list, u)

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
	w.WriteHeader(http.StatusCreated)
}
