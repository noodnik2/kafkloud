# syntax=docker/dockerfile:1

FROM golang:1.20.2-alpine3.17 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN apk add --no-progress --no-cache gcc musl-dev

ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

RUN go build -a -tags musl -o consumer .

FROM alpine:3.17

COPY --from=builder /app/consumer /
COPY docker.env .env

EXPOSE 8072

CMD ["/consumer"]