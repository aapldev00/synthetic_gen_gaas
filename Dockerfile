FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bin/generator-server ./cmd/server/main.go

FROM alpine:latest
COPY --from=builder /bin/generator-server /app/generator-server
EXPOSE 50051 9090
ENTRYPOINT ["/app/generator-server"]