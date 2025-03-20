FROM golang:1.24.1-alpine3.21 AS builder

RUN apk add --no-cache build-base

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o app ./cmd/app/main.go

FROM alpine:3.21.3

WORKDIR /app

COPY --from=builder /app/app .
COPY ./db/migrations ./migrations
ENV PG_MIGRATIONS_URI=file://migrations

ENTRYPOINT [ "./app" ]
