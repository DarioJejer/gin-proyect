# syntax=docker/dockerfile:1.7

FROM golang:1.25.6-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/gin-app .


FROM alpine:3.20

RUN apk add --no-cache ca-certificates \
  && addgroup -S app \
  && adduser -S app -G app

WORKDIR /app

COPY --from=builder /out/gin-app /app/gin-app

ENV GIN_MODE=release

EXPOSE 8080

USER app

ENTRYPOINT ["/app/gin-app"]

