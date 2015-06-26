FROM codeclimate/alpine-ruby:0.0.1

WORKDIR /usr/src/app
COPY bin/ /usr/src/app

RUN adduser -u 9000 -D app
USER app

CMD ["/usr/src/app/codeclimate-golint"]
