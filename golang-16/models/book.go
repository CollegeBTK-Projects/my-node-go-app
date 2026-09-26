package models

import "sync"

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   int    `json:"year"`
}

type BookUpdate struct {
	Title  *string `json:"title"`
	Author *string `json:"author"`
	Year   *int    `json:"year"`
}

type BookStore struct {
	mu     sync.Mutex
	books  []Book
	nextID int
}

func NewBookStore() *BookStore {
	return &BookStore{
		books: []Book{
			{ID: 1, Title: "Война и мир", Author: "Толстой", Year: 1869},
			{ID: 2, Title: "Преступление и наказание", Author: "Достоевский", Year: 1866},
			{ID: 3, Title: "Мастер и Маргарита", Author: "Булгаков", Year: 1967},
		},
		nextID: 4,
	}
}

func (s *BookStore) GetAll() []Book {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := make([]Book, len(s.books))
	copy(result, s.books)
	return result
}

func (s *BookStore) GetByID(id int) (Book, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, b := range s.books {
		if b.ID == id {
			return b, true
		}
	}
	return Book{}, false
}

func (s *BookStore) Add(title, author string, year int) Book {
	s.mu.Lock()
	defer s.mu.Unlock()
	book := Book{ID: s.nextID, Title: title, Author: author, Year: year}
	s.nextID++
	s.books = append(s.books, book)
	return book
}

func (s *BookStore) Update(id int, input BookUpdate) (Book, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.books {
		if s.books[i].ID == id {
			if input.Title != nil {
				s.books[i].Title = *input.Title
			}
			if input.Author != nil {
				s.books[i].Author = *input.Author
			}
			if input.Year != nil {
				s.books[i].Year = *input.Year
			}
			return s.books[i], true
		}
	}
	return Book{}, false
}

func (s *BookStore) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, b := range s.books {
		if b.ID == id {
			s.books = append(s.books[:i], s.books[i+1:]...)
			return true
		}
	}
	return false
}

func (s *BookStore) SearchByAuthor(author string) []Book {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := []Book{}
	for _, b := range s.books {
		if b.Author == author {
			result = append(result, b)
		}
	}
	return result
}
