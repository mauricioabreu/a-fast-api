FROM golang:1.20

WORKDIR /app

COPY go.mod go.sum /app/

RUN go mod download

CMD ["go", "run", "main.go"]
