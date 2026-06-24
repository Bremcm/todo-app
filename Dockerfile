FROM golang:1.26.4 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o todo-app ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/todo-app .
COPY --from=builder /app/configs ./configs

EXPOSE 8000

CMD ["./todo-app"]
