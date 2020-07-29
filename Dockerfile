FROM golang:latest AS builder
COPY . /app
RUN cd /app && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /blog

FROM alpine

COPY --from=builder /app/.env /.env
COPY --from=builder /blog /blog


ENTRYPOINT /blog