# Build stage
FROM golang:alpine AS builder
WORKDIR /usr/src/app
COPY . .
RUN go build -o job .

# Runtime stage
FROM alpine:latest
WORKDIR /usr/src/app
COPY --from=builder /usr/src/app/job .
ENTRYPOINT ["./job"]
