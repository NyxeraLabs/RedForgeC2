# Build stage
FROM golang:1.26-alpine AS builder
WORKDIR /app

COPY teamserver/go.mod teamserver/go.sum ./teamserver/
WORKDIR /app/teamserver
RUN go mod download

COPY teamserver ./
RUN go build -o /bin/teamserver ./cmd/teamserver

# Runtime stage
FROM alpine:3.19
RUN apk add --no-cache ca-certificates
COPY --from=builder /bin/teamserver /usr/local/bin/teamserver

EXPOSE 9080
ENTRYPOINT ["/usr/local/bin/teamserver"]
