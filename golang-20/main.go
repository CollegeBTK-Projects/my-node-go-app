package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"text/template"
	"time"
)

func LoggingMiddle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tm := time.Now()

		next.ServeHTTP(w, r)

		tmEnd := time.Since(tm).Microseconds()
		log.Printf("[%s] %dмс %s\n", r.Method, tmEnd, r.URL)
	}
}

func GroupHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Нету такого!", http.StatusNotFound)
		return
	}
	var storg bytes.Buffer

	tmpl, err := template.ParseFiles("components/group.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(&storg, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	storg.WriteTo(w)
}

func AboutHandler(w http.ResponseWriter, r *http.Request) {
	var storg bytes.Buffer

	tmpl, err := template.ParseFiles("components/about.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(&storg, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	storg.WriteTo(w)
}

func HeaderHandler(w http.ResponseWriter, r *http.Request) {
	js, err := json.Marshal(r.Header)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(js)
}

type CheckResponse struct {
	HeaderChecked string `json:"header_checked"`
	Exists        bool   `json:"exists"`
	CurrentValue  string `json:"current_value,omitempty"`
}

func SetHeadersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Powered-By", "Go")
	w.Header().Set("X-Group", "478")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("X-Temporary-Header", "Del")

	w.Header().Del("X-Temporary-Header")

	responseBody := map[string]string{
		"message": "Headers sets",
		"status":  "success",
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseBody)
}

func CheckHeadersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	targetHeader := "X-Group"
	w.Header().Set(targetHeader, "478")
	canonicalKey := http.CanonicalHeaderKey(targetHeader)
	_, exists := w.Header()[canonicalKey]
	var value string
	if exists {
		value = w.Header().Get(targetHeader)
	}
	result := CheckResponse{
		HeaderChecked: targetHeader,
		Exists:        exists,
		CurrentValue:  value,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func JsonHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Powered-By", "Go")
	w.Header().Set("X-Group-", "478")
	w.Header().Set("Cache-Control", "no-cache")
	response := map[string]string{"status": "success", "message": "Hello World"}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func HtmlHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("X-Powered-By", "Go")
	w.Header().Set("X-Group-", "478")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "<h1>HTML</h1>")
}

const max = 1 * 1024 * 1024

func EchoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, max)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		if err.Error() == "http: request body too large" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusRequestEntityTooLarge)
			json.NewEncoder(w).Encode(map[string]string{"error": "Max 1MB"})
			return
		}
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	log.Printf("[LOG] /echo | Размер тела: %d байт\n", len(body))

	contentType := strings.Split(r.Header.Get("Content-Type"), ";")[0]
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))

	switch contentType {
	case "application/json":
		var data map[string]interface{}
		if err := json.Unmarshal(body, &data); err != nil {
			http.Error(w, "JSON", http.StatusBadRequest)
			return
		}
		data["group"] = "478"

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)

	default:
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}
}

func FormHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, max)
	if err := r.ParseForm(); err != nil {
		if err.Error() == "http: request body too large" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusRequestEntityTooLarge)
			json.NewEncoder(w).Encode(map[string]string{"error": "Max 1MB"})
			return
		}
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}
	log.Printf("[LOG] /form | Успешно распарсено полей: %d\n", len(r.PostForm))
	parsedForm := make(map[string]string)
	for key, values := range r.PostForm {
		if len(values) > 0 {
			parsedForm[key] = values[0]
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(parsedForm)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", LoggingMiddle(GroupHandler))
	mux.HandleFunc("GET /about", LoggingMiddle(AboutHandler))
	mux.HandleFunc("GET /headers", LoggingMiddle(HeaderHandler))
	mux.HandleFunc("GET /headers/set", LoggingMiddle(SetHeadersHandler))
	mux.HandleFunc("GET /headers/check", LoggingMiddle(CheckHeadersHandler))
	mux.HandleFunc("GET /json", LoggingMiddle(JsonHandler))
	mux.HandleFunc("GET /html", LoggingMiddle(HtmlHandler))
	mux.HandleFunc("POST /echo", LoggingMiddle(EchoHandler))
	mux.HandleFunc("POST /form", LoggingMiddle(FormHandler))
	log.Fatal(http.ListenAndServe("localhost:3000", mux))
}
