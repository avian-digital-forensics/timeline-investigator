FROM golang:latest as builder
WORKDIR /tmp/app_build
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o console ./cmd/test-console/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /tmp/app_build/console ./console
CMD ./console --ti-url=http://34.102.155.104/api/