FROM alpine:3.2

ADD build/codeclimate-golint /usr/src/app/

RUN adduser -u 9000 -D app
USER app

CMD ["/usr/src/app/codeclimate-golint"]
