FROM golang:1.9.3-alpine3.7 as build

WORKDIR /usr/src/app

ENV LINT_VERSION=e14d9b0f1d332b1420c1ffa32562ad2dc84d645d

COPY codeclimate-golint.go /usr/src/app/codeclimate-golint.go
RUN apk add --no-cache git
RUN go get -d -t -v .
RUN (cd $GOPATH/src/github.com/golang/lint && git checkout ${LINT_VERSION} )
RUN go build -o codeclimate-golint .

COPY engine.json ./engine.json.template
RUN apk add --no-cache jq
RUN cat engine.json.template | jq '.version = .version + "/" + env.LINT_VERSION' > ./engine.json


FROM alpine:3.7

LABEL maintainer="Code Climate <hello@codeclimate.com>"

RUN adduser -u 9000 -D app

WORKDIR /usr/src/app

COPY --from=build /usr/src/app/engine.json /engine.json
COPY --from=build /usr/src/app/codeclimate-golint ./

USER app
WORKDIR /code
VOLUME /code

CMD ["/usr/src/app/codeclimate-golint"]
