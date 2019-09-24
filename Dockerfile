FROM golang:1.12-alpine as builder

RUN apk add dep git

RUN mkdir -p /go/src/github.com/bassman7689/honeybadger
WORKDIR /go/src/github.com/bassman7689/honeybadger

COPY . .

RUN dep ensure

RUN go build -o honeybadger cmd/smtpd/main.go
RUN go build -o webserver cmd/webserver/main.go

FROM alpine:latest

COPY --from=builder /go/src/github.com/bassman7689/honeybadger/public /static
COPY --from=builder /go/src/github.com/bassman7689/honeybadger/honeybadger /usr/local/bin/honeybadger
COPY --from=builder /go/src/github.com/bassman7689/honeybadger/webserver /usr/local/bin/webserver
CMD ["/usr/local/bin/honeybadger"]
