# base go image
FROM golang:1.21-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o gitlabServiceApp ./cmd

RUN chmod +x /app/gitlabServiceApp

# build a tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/gitlabServiceApp /app

CMD "/app/gitlabServiceApp"