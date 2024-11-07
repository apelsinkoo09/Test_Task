FROM golang:1.21.7

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

COPY configs/config.json /test_task/configs/config.json

RUN go build -o app ./cmd/server/main.go

CMD ["./app"]