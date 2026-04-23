FROM golang:1.25-alpine AS builder
WORKDIR /golang

RUN go install github.com/air-verse/air@latest

COPY ./go.mod ./go.sum* ./
RUN go mod download

EXPOSE 8080

CMD ["air"]
