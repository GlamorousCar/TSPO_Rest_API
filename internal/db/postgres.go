package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"tspo_server/model"
)

type BookRepository struct {
	db *sql.DB
}

func NewBookRepository(db *sql.DB) (*BookRepository, error) {
	return &BookRepository{db: db}, nil
}

func (r *BookRepository) GetBooks() ([]model.Book, error) {
	rows, err := r.db.Query("SELECT id, title, author FROM books")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []model.Book
	for rows.Next() {
		var book model.Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Author); err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}

func (r *BookRepository) GetBook(id string) (*model.Book, error) {
	var book model.Book
	err := r.db.QueryRow("SELECT id, title, author FROM books WHERE id = $1", id).
		Scan(&book.ID, &book.Title, &book.Author)
	if err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *BookRepository) CreateBook(book *model.Book) error {
	_, err := r.db.Exec("INSERT INTO books (id, title, author) VALUES ($1, $2, $3)",
		book.ID, book.Title, book.Author)
	return err
}

func (r *BookRepository) UpdateBook(book *model.Book) error {
	_, err := r.db.Exec("UPDATE books SET title = $1, author = $2 WHERE id = $3",
		book.Title, book.Author, book.ID)
	return err
}

func (r *BookRepository) DeleteBook(id string) error {
	_, err := r.db.Exec("DELETE FROM books WHERE id = $1", id)
	return err
}
