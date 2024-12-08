package app

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"strings"
	"tspo_server/internal/db"
	"tspo_server/model"
)

type Handler struct {
	repo   *db.BookRepository
	logger *slog.Logger
}

func NewHandler(repo *db.BookRepository, logger *slog.Logger) *Handler {
	return &Handler{
		repo:   repo,
		logger: logger,
	}
}

func (h *Handler) GetBooks(w http.ResponseWriter, r *http.Request) {
	defer func() {
		log.Println("error in GetBooks")
		recover()
	}()
	log.Println("GetBooks do")
	books, err := h.repo.GetBooks()
	log.Println("GetBooks after")
	if err != nil {
		h.logger.Error("failed to get books", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(books)
}

func (h *Handler) GetBook(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/books/")
	book, err := h.repo.GetBook(id)
	if err != nil {
		h.logger.Error("failed to get book", "error", err, "id", id)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(book)
}

func (h *Handler) CreateBook(w http.ResponseWriter, r *http.Request) {
	var book model.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		h.logger.Error("failed to decode request", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.repo.CreateBook(&book); err != nil {
		h.logger.Error("failed to create book", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

func (h *Handler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/books/")
	var book model.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		h.logger.Error("failed to decode request", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	book.ID = id
	if err := h.repo.UpdateBook(&book); err != nil {
		h.logger.Error("failed to update book", "error", err, "id", id)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(book)
}

func (h *Handler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/books/")
	if err := h.repo.DeleteBook(id); err != nil {
		h.logger.Error("failed to delete book", "error", err, "id", id)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
