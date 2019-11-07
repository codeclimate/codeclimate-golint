ARG BASE=1.13.1-alpine3.10
FROM golang:${BASE} as build

WORKDIR /usr/src/app

COPY engine.json ./engine.json.template
 RUN apk add --no-cache jq
 RUN export go_version=$(go version | cut -d ' ' -f 3) && \
     cat engine.json.template | jq '.version = .version + "/" + env.go_version' > ./engine.json

COPY codeclimate-golint.go go.mod go.sum ./
RUN apk add --no-cache git
RUN go build -o codeclimate-golint .

FROM golang:${BASE}
LABEL maintainer="Code Climate <hello@codeclimate.com>"

RUN adduser -u 9000 -D app

WORKDIR /usr/src/app

COPY --from=build /usr/src/app/engine.json /
COPY --from=build /usr/src/app/codeclimate-golint ./

USER app
VOLUME /code

CMD ["/usr/src/app/codeclimate-golint"]
