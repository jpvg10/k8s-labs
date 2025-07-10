FROM golang:alpine

ENV GIN_MODE=release
ENV PORT=8080

WORKDIR /usr/src/app

COPY . .

RUN go get ./
RUN go build .

EXPOSE $PORT

ENTRYPOINT ["./api"]
