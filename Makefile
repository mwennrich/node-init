GO111MODULE := on
DOCKER_TAG := $(or ${GITHUB_TAG_NAME}, latest)

all: node-init

.PHONY: node-init
node-init:
	CGO_ENABLED=0 go build -o bin/node-init -ldflags '-extldflags "-static"'
	strip bin/node-init

.PHONY: dockerimages
dockerimages:
	docker build -t metal-stack/node-init:${DOCKER_TAG} .

.PHONY: dockerpush
dockerpush:
	docker push metal-stack/node-init:${DOCKER_TAG}

.PHONY: clean
clean:
	rm -f bin/*

