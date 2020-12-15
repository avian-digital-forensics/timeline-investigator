FROM golang:latest as builder
WORKDIR /tmp/app_build
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api ./cmd/main/main.go
