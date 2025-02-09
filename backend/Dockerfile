FROM golang:alpine AS builder
ENV CGO_ENABLED=0
ENV GOOS=linux
RUN apk update --no-cache && apk add --no-cache git
WORKDIR /build
ADD go.mod .
ADD go.sum .
RUN go mod download
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
COPY . .
RUN go build -ldflags="-s -w" -o /app/main ./cmd/main.go

FROM alpine
RUN apk update --no-cache
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/main /app/main
COPY --from=builder /build/migrations /app/migrations
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY docker-entrypoint.sh .
RUN chmod +x /app/docker-entrypoint.sh
ENTRYPOINT ["/app/docker-entrypoint.sh"]
