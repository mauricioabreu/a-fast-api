FROM golang:1.20 as builder

WORKDIR /app

COPY go.mod go.sum /app/

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -v -o main .

FROM alpine:3.18.3

EXPOSE 80

COPY --from=builder /app/main .

CMD ["/main"]
