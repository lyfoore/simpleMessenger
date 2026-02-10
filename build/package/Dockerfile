FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN go build -o main cmd/simpleMessenger/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main ./
CMD ["./main"]