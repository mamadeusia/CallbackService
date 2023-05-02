# FROM golang:alpine AS builder
FROM golang:1.20-bullseye AS base
WORKDIR /go/src/app
COPY go.mod go.sum ./

COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
    go mod download -x
RUN go build 

FROM debian:bullseye AS dist
RUN apt update && apt install -y ca-certificates
WORKDIR /app
RUN groupadd -g 1001 -r micro && \
        useradd -u 1001 -r -s /bin/false -d /app -g micro micro && \
        chown -R micro:micro /app
USER micro:micro
COPY --from=base --chown=micro:micro /go/src/app/CallbackService /app
CMD /app/connector-callback-service
