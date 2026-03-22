FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o xray-exporter .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/xray-exporter .

ARG XRAY_ENDPOINT=localhost:11111
ARG PORT=9595

CMD ["./xray-exporter", "-xray-endpoint", "${XRAY_ENDPOINT}", "-port", "${PORT}"]