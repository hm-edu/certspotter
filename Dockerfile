# syntax=docker/dockerfile:1

FROM golang:1.26.5
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o certspotter ./cmd/certspotter/

FROM alpine:3.24.1@sha256:28bd5fe8b56d1bd048e5babf5b10710ebe0bae67db86916198a6eec434943f8b
COPY --from=0 /app/certspotter /usr/local/bin/certspotter