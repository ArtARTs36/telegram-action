# syntax=docker/dockerfile:1

FROM golang:1.23.3-alpine AS builder

ARG APP_VERSION="undefined"
ARG BUILD_TIME="undefined"

WORKDIR /go/src/github.com/artarts36/telegram-action

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOOS=linux go build -ldflags="-s -w -X 'main.Version=${APP_VERSION}' -X 'main.BuildDate=${BUILD_TIME}'" -o /go/bin/telegram-action /go/src/github.com/artarts36/telegram-action/cmd/main.go

######################################################

FROM alpine

RUN apk add tini

COPY --from=builder /go/bin/telegram-action /go/bin/telegram-action

# https://github.com/opencontainers/image-spec/blob/main/annotations.md
LABEL org.opencontainers.image.title="telegram-action"
LABEL org.opencontainers.image.description="send telegram message"
LABEL org.opencontainers.image.url="https://github.com/artarts36/telegram-action"
LABEL org.opencontainers.image.source="https://github.com/artarts36/telegram-action"
LABEL org.opencontainers.image.vendor="ArtARTs36"
LABEL org.opencontainers.image.version="$APP_VERSION"
LABEL org.opencontainers.image.created="$BUILD_TIME"
LABEL org.opencontainers.image.licenses="MIT"

COPY docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x ./docker-entrypoint.sh

ENTRYPOINT ["/docker-entrypoint.sh"]
