FROM golang:latest as builder
WORKDIR /tmp/app_build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api ./cmd/main/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /tmp/app_build/api ./api
COPY --from=builder /tmp/app_build/.local/config.local.yml /configs/config.yml
CMD ./api --cfg=/configs/config.yml