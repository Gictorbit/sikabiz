# syntax=docker/dockerfile:1.4
ARG GO_VERSION="1.22"
ARG GOPROXYURL="https://goproxy.io,direct"
ARG COMPRESS="true"
ARG GITHUB_TOKEN
ARG GITHUB_ORIGIN=""

FROM golang:${GO_VERSION}-alpine AS builder
# install packages
RUN sed -i 's#dl-cdn.alpinelinux.org#alpine.global.ssl.fastly.net#g' /etc/apk/repositories
RUN apk --no-cache add --update ca-certificates tzdata upx git

# config git
ARG GITHUB_TOKEN
#ENV GOPRIVATE="${COMPANY_HOST}/*"
RUN git config --global url."https://${GITHUB_TOKEN}:x-oauth-basic@github.com".insteadOf "https://github.com"

# copy source code
WORKDIR /build
COPY . .

# Get all of our dependencies
ARG GOPROXYURL
RUN --mount=type=cache,mode=0755,target=/go/pkg/mod GOPROXY="${GOPROXYURL}" go mod download -x
# compile project
RUN --mount=type=cache,mode=0755,target=/go/pkg/mod CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags "-s -w" -a -installsuffix cgo -o ./bin/usersrv ./cmd/...

ARG COMPRESS
RUN mkdir -p /final && \
    if [ "$COMPRESS" = "true" ] ;then upx --best --lzma -o /final/usersrv ./bin/usersrv ;else cp ./bin/usersrv /final; fi

FROM alpine:latest AS final

WORKDIR /production
COPY --from=builder /final .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

ARG GITHUB_ORIGIN
LABEL org.opencontainers.image.source="${GITHUB_ORIGIN}"

ENTRYPOINT ["./usersrv"]
