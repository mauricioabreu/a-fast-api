FROM golang:1.20

WORKDIR /app

RUN go install github.com/cosmtrek/air@latest

COPY go.mod go.sum /app/

RUN go mod download

CMD ["air", "-c", ".air.toml"]
