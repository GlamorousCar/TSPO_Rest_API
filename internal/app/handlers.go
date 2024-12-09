package app

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"
	"tspo_server/internal/db"
	"tspo_server/internal/errors"
	"tspo_server/internal/query"
	"tspo_server/model"
)

type Handler struct {
	repo   *db.BookRepository
	logger *slog.Logger
}

type Response struct {
	Data       interface{}      `json:"data,omitempty"`
	Pagination *Pagination      `json:"pagination,omitempty"`
	Error      *errors.APIError `json:"error,omitempty"`
}

type Pagination struct {
	CurrentPage  int `json:"current_page"`
	PageSize     int `json:"page_size"`
	TotalPages   int `json:"total_pages"`
	TotalRecords int `json:"total_records"`
}

func NewHandler(repo *db.BookRepository, logger *slog.Logger) *Handler {
	return &Handler{
		repo:   repo,
		logger: logger,
	}
}

func (h *Handler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) writeError(w http.ResponseWriter, err error) {
	var status int
	var message string

	switch {
	case err == errors.ErrNotFound:
		status = http.StatusNotFound
		message = "Resource not found"
	case err == errors.ErrInvalidInput:
		status = http.StatusBadRequest
		message = "Invalid input"
	case err == errors.ErrTimeout:
		status = http.StatusGatewayTimeout
		message = "Operation timed out"
	default:
		status = http.StatusInternalServerError
		message = "Internal server error"
	}

	h.writeJSON(w, status, Response{
		Error: errors.NewAPIError(status, message),
	})
}

func (h *Handler) GetBooks(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	params := query.NewParams(r)
	books, total, err := h.repo.GetBooks(ctx, params)
	if err != nil {
		h.logger.Error("failed to get books", "error", err)
		h.writeError(w, err)
		return
	}

	totalPages := (total + params.PageSize - 1) / params.PageSize
	h.writeJSON(w, http.StatusOK, Response{
		Data: books,
		Pagination: &Pagination{
			CurrentPage:  params.Page,
			PageSize:     params.PageSize,
			TotalPages:   totalPages,
			TotalRecords: total,
		},
	})
}

func (h *Handler) GetBook(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	id := strings.TrimPrefix(r.URL.Path, "/books/")
	book, err := h.repo.GetBook(ctx, id)
	if err != nil {
		h.logger.Error("failed to get book", "error", err, "id", id)
		h.writeError(w, err)
		return
	}

	h.writeJSON(w, http.StatusOK, Response{Data: book})
}

func (h *Handler) CreateBook(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var book model.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		h.logger.Error("failed to decode request", "error", err)
		h.writeError(w, errors.ErrInvalidInput)
		return
	}

	if err := h.repo.CreateBook(ctx, &book); err != nil {
		h.logger.Error("failed to create book", "error", err)
		h.writeError(w, err)
		return
	}

	h.writeJSON(w, http.StatusCreated, Response{Data: book})
}

func (h *Handler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	id := strings.TrimPrefix(r.URL.Path, "/books/")
	var book model.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		h.logger.Error("failed to decode request", "error", err)
		h.writeError(w, errors.ErrInvalidInput)
		return
	}

	book.ID = id
	if err := h.repo.UpdateBook(ctx, &book); err != nil {
		h.logger.Error("failed to update book", "error", err, "id", id)
		h.writeError(w, err)
		return
	}

	h.writeJSON(w, http.StatusOK, Response{Data: book})
}

func (h *Handler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	id := strings.TrimPrefix(r.URL.Path, "/books/")
	if err := h.repo.DeleteBook(ctx, id); err != nil {
		h.logger.Error("failed to delete book", "error", err, "id", id)
		h.writeError(w, err)
		return
	}

	h.writeJSON(w, http.StatusNoContent, nil)
}
