# backend/Dockerfile.backend
FROM golang:1.21-alpine as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o autosysadmin ./cmd/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/autosysadmin .
COPY --from=builder /app/migrations ./migrations
EXPOSE 8080
CMD ["./autosysadmin"]