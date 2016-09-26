FROM alpine:3.2

MAINTAINER Code Climate <hello@codeclimate.com>

WORKDIR /code

COPY build/codeclimate-golint /usr/src/app/

VOLUME /code

RUN adduser -u 9000 -D app
USER app

CMD ["/usr/src/app/codeclimate-golint"]
