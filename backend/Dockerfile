FROM golang:1.23.4 AS builder
WORKDIR /app

COPY . .
RUN go mod tidy
RUN go build -o go-server .

FROM debian:latest
WORKDIR /root/
COPY --from=builder /app/go-server .
EXPOSE 8080
CMD ["./go-server"]