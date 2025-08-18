# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build application binary
RUN CGO_ENABLED=0 GOOS=linux go build -o linko ./cmd/linko

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

COPY --from=builder /app/linko .

VOLUME ["/etc/linko"]

EXPOSE 8080

ENTRYPOINT ["./linko"]
CMD ["-config.file=/etc/linko/config.yaml"]