FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /bin/diploma cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /bin/diploma .

CMD ["./diploma"]
