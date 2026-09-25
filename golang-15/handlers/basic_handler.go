package handlers

import (
	"net/http"
	"text/template"
	"time"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("components/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tm := time.Now().Format("02-01-2006 15:04:03")

	if err := tmpl.Execute(w, tm); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
