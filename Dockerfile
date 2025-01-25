FROM golang:1.20-alpine

WORKDIR /app

ENV GO111MODULE=on

COPY go.mod go.sum ./
RUN go mod download

COPY ./api ./api

WORKDIR /app/api

RUN go build -o main main.go

CMD ["./main"]