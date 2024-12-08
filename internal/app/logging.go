package app

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

// https://stackoverflow.com/questions/53272536/how-do-i-get-response-statuscode-in-golang-middleware
type CustomResponseWriter struct {
	responseWriter http.ResponseWriter
	StatusCode     int
	body           *bytes.Buffer // Буфер для записи тела ответа
}

func ExtendResponseWriter(w http.ResponseWriter) *CustomResponseWriter {
	return &CustomResponseWriter{
		responseWriter: w,
		StatusCode:     0,
		body:           &bytes.Buffer{},
	}
}

func (w *CustomResponseWriter) Write(b []byte) (int, error) {
	// Сохраняем данные в буфер
	w.body.Write(b)
	// Записываем данные в реальный ResponseWriter
	return w.responseWriter.Write(b)
}

func (w *CustomResponseWriter) Header() http.Header {
	return w.responseWriter.Header()
}

func (w *CustomResponseWriter) WriteHeader(statusCode int) {
	// Устанавливаем статус ответа
	w.StatusCode = statusCode
	w.responseWriter.WriteHeader(statusCode)
}

func (w *CustomResponseWriter) Done() {
	// Если WriteHeader не был вызван, устанавливаем статус 200 OK
	if w.StatusCode == 0 {
		w.StatusCode = http.StatusOK
	}
}

func HandlerLogging(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var agent string
			if a := r.UserAgent(); a != "" {
				agent = fmt.Sprintf("using %s", a)
			}
			user := r.Header.Get("X-Remote-User")

			if strings.HasPrefix(r.URL.Path, "/webdav") {
				logger.Debug("Request received with /webdav prefix")
				r.URL.Path = strings.Replace(r.URL.Path, "/webdav", "/", 1)
				//r.URL.Path = strings.TrimPrefix(r.URL.Path, "/webdav")
			}

			// Дополнительная логика для уровня Debug
			if logger.Enabled(context.Background(), slog.LevelDebug) {

				body, err := io.ReadAll(r.Body)

				defer r.Body.Close() // Закрываем тело запроса после чтения

				// Преобразуем тело в строку
				bodyStr := string(body)

				// Восстановление тела запроса (если нужно передать дальше)
				r.Body = io.NopCloser(bytes.NewReader(body))

				if err == nil {
					bodyStr = string(bodyStr)
				} else {
					// Логируем ошибку чтения тела
					bodyStr = err.Error()
				}

				headerMap := make(map[string]interface{})
				for name, values := range r.Header {
					if len(values) > 1 {
						headerMap[name] = values
					} else {
						headerMap[name] = values[0]
					}
				}
				logger.LogAttrs(
					context.Background(),
					slog.LevelDebug,
					"Request received (detailed information)",
					slog.String("request_body", bodyStr),
					slog.String("method", r.Method),
					slog.String("path", r.URL.Path),
					slog.String("user", user),
					slog.String("remote_addr", r.RemoteAddr),
					slog.Any("headers", headerMap),
				)

			} else {
				logger.LogAttrs(
					context.Background(),
					slog.LevelInfo,
					"Request received",
					slog.String("method", r.Method),
					slog.String("path", r.URL.Path),
					slog.String("user", user),
					slog.String("remote_addr", r.RemoteAddr),
					slog.String("agent", agent),
				)
			}

			//logger.Info(fmt.Sprintf("%s request for %s by %s received from %s %s", r.Method, r.URL.Path, user, r.RemoteAddr, agent))

			//logger.Debug("\n", r.Method, r.URL.Path, user, "----req.Body----\n", "\n-------")
			//defer r.Body.Close()
			ew := ExtendResponseWriter(w)

			startTime := time.Now()

			next.ServeHTTP(ew, r)
			duration := time.Since(startTime)

			if ew.StatusCode >= 400 {
				logger.LogAttrs(
					context.Background(),
					slog.LevelInfo,
					"Response sent",
					slog.String("method", r.Method),
					slog.String("path", r.URL.Path),
					slog.String("user", user),
					slog.String("remote_addr", r.RemoteAddr),
					slog.Int("status_code", ew.StatusCode),
					slog.String("response_body", ew.body.String()),
					slog.Duration("duration", duration),
				)
			} else {
				logger.LogAttrs(
					context.Background(),
					slog.LevelInfo,
					"Response sent",
					slog.String("method", r.Method),
					slog.String("path", r.URL.Path),
					slog.String("user", user),
					slog.String("remote_addr", r.RemoteAddr),
					slog.Int("status_code", ew.StatusCode),
					slog.String("duration", duration.String()),
				)
			}

			ew.Done()

		})
	}
}
