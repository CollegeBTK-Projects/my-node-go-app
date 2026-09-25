package handlers

import (
	"html/template"
	"net/http"
	"time"
)

var templates = template.Must(template.ParseGlob("components/*.html"))

func Home(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"Date": time.Now().Format("02.01.2006 15:04:05"),
	}
	templates.ExecuteTemplate(w, "home.html", data)
}

func About(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "about.html", nil)
}

