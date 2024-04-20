FROM golang:alpine3.19 AS builder

LABEL stage=appbuilder

ENV CGO_ENABLED 0

ENV GOOS linux

RUN apk update --no-cache

WORKDIR /builder

COPY go.mod .

COPY go.sum .

RUN go mod download

COPY . .

RUN go build -C ./cmd/app -o ../../build/app/main -ldflags="-s -w"

FROM alpine:3.19.1

LABEL stage=apprunner

RUN apk update --no-cache && apk add --no-cache ca-certificates && apk add --no-cache --upgrade bash

WORKDIR /build

COPY . .

COPY --from=builder /builder/build/app/main /build

CMD ["./main"]
