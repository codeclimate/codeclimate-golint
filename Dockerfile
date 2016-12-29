FROM alpine:3.4

MAINTAINER Code Climate <hello@codeclimate.com>

WORKDIR /usr/src/app
COPY codeclimate-golint.go /usr/src/app

RUN apk --update add go git && \
  export GOPATH=/tmp/go GOBIN=/usr/local/bin && \
  go get -d . && \
  go install codeclimate-golint.go && \
  apk del go git && \
  rm -rf "$GOPATH" rm /var/cache/apk/*

WORKDIR /code
VOLUME /code

RUN adduser -u 9000 -D app
USER app

CMD ["/usr/local/bin/codeclimate-golint"]
