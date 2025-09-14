FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY . .

# No need for `go mod init`, assuming go.mod is already in the repo
RUN go mod tidy
RUN go build -o smackapi main.go

FROM alpine:3.17
WORKDIR /app
COPY --from=builder /app/smackapi .

ENV VERSION=v1
EXPOSE 8080

CMD ["./smackapi"]
