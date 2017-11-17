FROM alpine:edge

LABEL maintainer="Code Climate <hello@codeclimate.com>"

RUN adduser -u 9000 -D app

WORKDIR /usr/src/app

COPY engine.json codeclimate-golint.go ./

RUN apk add --no-cache --virtual .dev-deps musl-dev go git jq && \
  export GOPATH=/tmp/go GOBIN=/usr/local/bin && \
  go get -d -t -v . && \
  export LIBRARY_PATH=/usr/lib32:$LIBRARY_PATH && \
  go install codeclimate-golint.go && \
  export golint_version=$(cd "${GOPATH}/src/github.com/golang/lint/" && git rev-parse HEAD 2>/dev/null) && \
  cat engine.json | jq '.version = .version + "/" + env.golint_version' > /engine.json && \
  apk del .dev-deps && \
  rm -rf "$GOPATH"

USER app
WORKDIR /code
VOLUME /code

CMD ["/usr/local/bin/codeclimate-golint"]
