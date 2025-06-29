# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o minimaldo-backend .

# Run stage
FROM alpine:latest
WORKDIR /root/

COPY --from=builder /app/minimaldo-backend .
EXPOSE 8080
CMD ["./minimaldo-backend"]
