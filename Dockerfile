# syntax=docker/dockerfile:1

FROM golang:1.23.5
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o certspotter ./cmd/certspotter/

FROM alpine:3.21.1@sha256:b97e2a89d0b9e4011bb88c02ddf01c544b8c781acf1f4d559e7c8f12f1047ac3
COPY --from=0 /app/certspotter /usr/local/bin/certspotter