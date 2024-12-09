package db

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"strings"
	"tspo_server/internal/errors"
	"tspo_server/internal/query"
	"tspo_server/model"
)

type BookRepository struct {
	db *sql.DB
}

func NewBookRepository(db *sql.DB) (*BookRepository, error) {
	return &BookRepository{db: db}, nil
}

func (r *BookRepository) GetBooks(ctx context.Context, params *query.Params) ([]model.Book, int, error) {
	// Build the query with filters
	whereClause := []string{}
	args := []interface{}{}
	argCount := 1

	for key, value := range params.Filter {
		whereClause = append(whereClause, fmt.Sprintf("%s ILIKE $%d", key, argCount))
		args = append(args, "%"+value+"%")
		argCount++
	}

	countQuery := "SELECT COUNT(*) FROM books"
	if len(whereClause) > 0 {
		countQuery += " WHERE " + strings.Join(whereClause, " AND ")
	}

	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("%w: %v", errors.ErrDatabaseOperation, err)
	}

	query := "SELECT id, title, author FROM books"
	if len(whereClause) > 0 {
		query += " WHERE " + strings.Join(whereClause, " AND ")
	}

	query += fmt.Sprintf(" ORDER BY %s %s", params.Sort, params.Order)

	offset := (params.Page - 1) * params.PageSize
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, params.PageSize, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("%w: %v", errors.ErrDatabaseOperation, err)
	}
	defer rows.Close()

	var books []model.Book
	for rows.Next() {
		var book model.Book
		if err = rows.Scan(&book.ID, &book.Title, &book.Author); err != nil {
			return nil, 0, fmt.Errorf("%w: %v", errors.ErrDatabaseOperation, err)
		}
		books = append(books, book)
	}

	return books, total, nil
}

func (r *BookRepository) GetBook(ctx context.Context, id string) (*model.Book, error) {
	var book model.Book
	err := r.db.QueryRowContext(ctx, "SELECT id, title, author FROM books WHERE id = $1", id).
		Scan(&book.ID, &book.Title, &book.Author)

	if err == sql.ErrNoRows {
		return nil, errors.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrDatabaseOperation, err)
	}

	return &book, nil
}

func (r *BookRepository) CreateBook(ctx context.Context, book *model.Book) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO books (id, title, author) VALUES ($1, $2, $3)",
		book.ID, book.Title, book.Author)

	if err != nil {
		return fmt.Errorf("%w: %v", errors.ErrDatabaseOperation, err)
	}
	return nil
}

func (r *BookRepository) UpdateBook(ctx context.Context, book *model.Book) error {
	result, err := r.db.ExecContext(ctx,
		"UPDATE books SET title = $1, author = $2 WHERE id = $3",
		book.Title, book.Author, book.ID)

	if err != nil {
		return fmt.Errorf("%w: %v", errors.ErrDatabaseOperation, err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v", errors.ErrDatabaseOperation, err)
	}
	if rows == 0 {
		return errors.ErrNotFound
	}

	return nil
}

func (r *BookRepository) DeleteBook(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM books WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("%w: %v", errors.ErrDatabaseOperation, err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v", errors.ErrDatabaseOperation, err)
	}
	if rows == 0 {
		return errors.ErrNotFound
	}

	return nil
}
