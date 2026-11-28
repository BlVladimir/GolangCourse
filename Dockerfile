FROM golang:1.26.0 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/catalog-service \
    ./internal/cmd/app

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata && \
    adduser -D -H -u 10001 appuser

WORKDIR /app

COPY --from=builder /out/catalog-service /app/catalog-service

USER appuser
EXPOSE 8080 9090

ENTRYPOINT ["/app/catalog-service"]