# syntax=docker/dockerfile:1

FROM golang:1.22.3
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o certspotter ./cmd/certspotter/

FROM alpine:3.19.1@sha256:c5b1261d6d3e43071626931fc004f70149baeba2c8ec672bd4f27761f8e1ad6b
COPY --from=0 /app/certspotter /usr/local/bin/certspotter