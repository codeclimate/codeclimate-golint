.PHONY: update image release

IMAGE_NAME ?= codeclimate/codeclimate-golint
RELEASE_REGISTRY ?= codeclimate

ifndef RELEASE_TAG
override RELEASE_TAG = latest
endif

image:
	docker build --rm -t $(IMAGE_NAME) .

release:
	docker tag $(IMAGE_NAME) $(RELEASE_REGISTRY)/codeclimate-golint:$(RELEASE_TAG)
	docker push $(RELEASE_REGISTRY)/codeclimate-golint:$(RELEASE_TAG)
