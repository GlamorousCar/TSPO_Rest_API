FROM golang:1.23-alpine

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/tspo_server

EXPOSE 8080

CMD ["./server"]