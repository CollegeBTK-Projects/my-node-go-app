package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/CollegeBTK-Projects/my-node-go-app/middleware"
	"github.com/CollegeBTK-Projects/my-node-go-app/models"
)

type BookHandler struct {
	store *models.BookStore
}

func NewBookHandler(store *models.BookStore) *BookHandler {
	return &BookHandler{store: store}
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *BookHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.store.GetAll())
}

func (h *BookHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		panic(&middleware.AppError{Status: http.StatusBadRequest, Message: "Неверный ID книги"})
	}
	book, found := h.store.GetByID(id)
	if !found {
		panic(&middleware.AppError{Status: http.StatusNotFound, Message: "Книга не найдена"})
	}
	writeJSON(w, http.StatusOK, book)
}

func (h *BookHandler) Search(w http.ResponseWriter, r *http.Request) {
	author := r.URL.Query().Get("author")
	if author == "" {
		panic(&middleware.AppError{Status: http.StatusBadRequest, Message: "Укажите параметр author"})
	}
	writeJSON(w, http.StatusOK, h.store.SearchByAuthor(author))
}

func (h *BookHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title  string `json:"title"`
		Author string `json:"author"`
		Year   int    `json:"year"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		panic(&middleware.AppError{Status: http.StatusBadRequest, Message: "Невалидный JSON"})
	}
	if input.Title == "" || input.Author == "" {
		panic(&middleware.AppError{Status: http.StatusBadRequest, Message: "Поля title и author обязательны"})
	}
	if input.Year <= 0 || input.Year > 2100 {
		panic(&middleware.AppError{Status: http.StatusBadRequest, Message: "Некорректный год издания"})
	}
	book := h.store.Add(input.Title, input.Author, input.Year)
	writeJSON(w, http.StatusCreated, book) 
}

func (h *BookHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		panic(&middleware.AppError{Status: http.StatusBadRequest, Message: "Неверный ID книги"})
	}
	var input models.BookUpdate
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		panic(&middleware.AppError{Status: http.StatusBadRequest, Message: "Невалидный JSON"})
	}
	book, found := h.store.Update(id, input)
	if !found {
		panic(&middleware.AppError{Status: http.StatusNotFound, Message: "Книга не найдена"})
	}
	writeJSON(w, http.StatusOK, book)
}

func (h *BookHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		panic(&middleware.AppError{Status: http.StatusBadRequest, Message: "Неверный ID книги"})
	}
	if ok := h.store.Delete(id); !ok {
		panic(&middleware.AppError{Status: http.StatusNotFound, Message: "Книга не найдена"})
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Книга успешно удалена"})
}
