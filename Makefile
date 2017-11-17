.PHONY: update image

IMAGE_NAME ?= codeclimate/codeclimate-golint

update:
	docker run \
	  --rm --interactive \
	  -v $(PWD)/engine.json:/engine.json \
	  -v $(PWD)/bin/update:/usr/local/bin/update \
	  alpine:edge \
	  /usr/local/bin/update

image: update
	docker build --rm -t $(IMAGE_NAME) .
