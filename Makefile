.PHONY: update image

IMAGE_NAME ?= codeclimate/codeclimate-golint

image:
	docker build --rm -t $(IMAGE_NAME) .
