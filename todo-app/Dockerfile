FROM golang:alpine

ENV GIN_MODE=release

WORKDIR /usr/src/app

COPY . .

RUN go get ./
RUN go build .

ENTRYPOINT ["./api"]
