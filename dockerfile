FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o bot .

FROM alpine:3.19

WORKDIR /app
COPY --from=builder /app/bot .
ENTRYPOINT [ "./bot" ]
