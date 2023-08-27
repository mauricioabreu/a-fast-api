# builder
FROM golang:1.21.0-alpine3.18 as builder

WORKDIR /app

COPY go.mod go.sum /app/

RUN go mod download

RUN apk add dumb-init

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -v -o main .

# deploy
FROM alpine:3.18.3

EXPOSE 80

COPY --from=builder /usr/bin/dumb-init /usr/bin/dumb-init
COPY --from=builder /app/main .

ENTRYPOINT ["/usr/bin/dumb-init", "--"]

CMD ["/main"]

STOPSIGNAL SIGINT
