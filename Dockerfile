FROM alpine:edge

LABEL maintainer="Code Climate <hello@codeclimate.com>"

RUN adduser -u 9000 -D app

WORKDIR /usr/src/app

COPY engine.json /engine.json
COPY codeclimate-golint.go /usr/src/app/codeclimate-golint.go

RUN apk add --no-cache --virtual .dev-deps musl-dev go git && \
  export GOPATH=/tmp/go GOBIN=/usr/local/bin && \
  go get -d -t -v . && \
  export LIBRARY_PATH=/usr/lib32:$LIBRARY_PATH && \
  go install codeclimate-golint.go && \
  apk del .dev-deps && \
  rm -rf "$GOPATH"

USER app
WORKDIR /code
VOLUME /code

CMD ["/usr/local/bin/codeclimate-golint"]
