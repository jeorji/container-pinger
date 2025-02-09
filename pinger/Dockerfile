FROM golang:alpine AS builder
ENV CGO_ENABLED=0
ENV GOOS=linux
RUN apk update --no-cache && apk add --no-cache git
WORKDIR /build
ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
RUN go build -ldflags="-s -w" -o /app/main ./cmd/main.go

FROM alpine
RUN apk update --no-cache
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/main /app/main
RUN chmod +x /app/main
ENTRYPOINT ["/app/main"]
