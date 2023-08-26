# builder
FROM golang:1.20-buster as builder

WORKDIR /app

COPY go.mod go.sum /app/

RUN go mod download

RUN apt-get update && apt-get install -y dumb-init

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -v -o main .

# deploy
FROM gcr.io/distroless/base-debian11

EXPOSE 80

COPY --from=builder /usr/bin/dumb-init /usr/bin/dumb-init
COPY --from=builder /app/main .

USER nonroot:nonroot

ENTRYPOINT ["/usr/bin/dumb-init", "--"]

CMD ["/main"]
