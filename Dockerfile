FROM golang:1.20-alpine

WORKDIR /app

ENV GO111MODULE=on

COPY go.mod go.sum ./
RUN go mod download

COPY ./api ./api
COPY .env ./.env

WORKDIR /app/api

RUN go build -o main main.go

EXPOSE 8000

CMD ["./main"]