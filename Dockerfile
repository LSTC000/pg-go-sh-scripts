FROM golang:alpine3.19 AS builder

LABEL stage=appbuilder

ENV CGO_ENABLED 0

ENV GOOS linux

RUN apk update --no-cache

WORKDIR /builder

ADD go.mod .

ADD go.sum .

RUN go mod download

COPY . .

RUN go build -ldflags="-s -w" -C ./cmd/app -o ../../build/app/main

FROM alpine:3.19.1

LABEL stage=apprunner

RUN apk update --no-cache && apk add --no-cache ca-certificates

WORKDIR /build

COPY --from=builder /builder/build/app/main /build

CMD ["./main"]
