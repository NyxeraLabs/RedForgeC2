# Build stage
FROM golang:1.26-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /bin/teamserver ./cmd/teamserver

# Runtime stage
FROM alpine:3.19
RUN apk add --no-cache ca-certificates
COPY --from=builder /bin/teamserver /usr/local/bin/teamserver

EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/teamserver"]
