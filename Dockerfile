FROM golang:1.21 AS builder

RUN mkdir /app
ADD . /app
WORKDIR /app

COPY ./internal/postgres/migrations /app/migrations

RUN CGO_ENABLED=0 GOOS=linux go build -o app cmd/*.go

FROM alpine:latest AS production
COPY --from=builder /app .
CMD ["./app"]
