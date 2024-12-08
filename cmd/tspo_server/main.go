package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/mdobak/go-xerrors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"tspo_server/internal/app"
	"tspo_server/internal/config"
	"tspo_server/internal/db"
	"tspo_server/pkg/logger"
)

func main() {
	c := config.Configuration{}
	c.Construct()

	// инициализация логгирования
	logger := logger.CreateLogger(c.LogLevel)

	database, err := NewDB(&c)
	if err != nil {
		xerr := xerrors.New(err)
		logger.LogAttrs(context.Background(), slog.LevelError, "Failed to connect to database, check env variables or connections", slog.Any("error", xerr))
		return
	}
	defer database.Close()

	repo, err := db.NewBookRepository(database)

	// Initialize database

	handler := app.NewHandler(repo, logger)

	// Setup router
	mux := http.NewServeMux()

	// Apply JWT middleware to all book routes
	//jwtMiddleware := auth.NewJWTMiddleware(c.JWTSecret)

	// Register routes
	mux.HandleFunc("GET /books", handler.GetBooks)
	mux.HandleFunc("GET /books/{id}", handler.GetBook)
	mux.HandleFunc("POST /books", handler.CreateBook)
	mux.HandleFunc("PUT /books/{id}", handler.UpdateBook)
	mux.HandleFunc("DELETE /books/{id}", handler.DeleteBook)

	// Start server

	srv := &http.Server{
		Addr: os.Getenv("API_SERVER_ADDR"),
		//Handler: tracing(nextRequestID)(logging(loggerNew)(mux)),
		Handler: app.HandlerLogging(logger)(mux),
	}

	logger.Info("Starting server on", slog.String("server addr", srv.Addr))

	if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		xerr := xerrors.New(err)
		logger.LogAttrs(context.Background(), slog.LevelError, "Server failed to start", slog.Any("error", xerr))
		//TODO exit?
		os.Exit(1)
	}
}

func NewDB(c *config.Configuration) (*sql.DB, error) {
	var db *sql.DB
	var err error

	if c.DBFlavor == "postgres" {
		psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DBNAME"))

		log.Println("===", psqlInfo)
		db, err = sql.Open("postgres", psqlInfo)

	} else if c.DBFlavor == "sqlite3" {
		return nil, errors.New("sqlite3 not supported")

	}
	if err != nil {
		return nil, err
	}
	// Проверка подключения
	if err = db.Ping(); err != nil {
		return nil, err
	}
	log.Println("Connected to database successfully")
	return db, nil
}
