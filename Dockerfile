FROM golang:1.12-alpine as builder

RUN mkdir -p /go/src/github.com/bassman7689/honeybadger
WORKDIR /go/src/github.com/bassman7689/honeybadger

COPY . .

RUN go build -o honeybadger main.go

FROM alpine:latest

COPY --from=builder /go/src/github.com/bassman7689/honeybadger/honeybadger /usr/local/bin/honeybadger
CMD ["/usr/local/bin/honeybadger"]
