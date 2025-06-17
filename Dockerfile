# ----- BUILD STAGE -----
FROM golang:1.23.6 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -o main

# ----- RUN STAGE -----
FROM alpine:latest

WORKDIR /app

# Correct path: no "app/" prefix here — just absolute path from build container
COPY --from=builder /app/main .
COPY --from=builder /app/static ./static
COPY --from=builder /app/.env .env

RUN chmod +x main

EXPOSE 8080

CMD ["./main"]
