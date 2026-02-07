# Build stage
FROM golang:alpine AS builder
WORKDIR /usr/src/app
COPY . .
RUN go get ./
RUN go build -o broadcaster .

# Runtime stage
FROM alpine:latest
WORKDIR /usr/src/app
COPY --from=builder /usr/src/app/broadcaster .
ENTRYPOINT ["./broadcaster"]
