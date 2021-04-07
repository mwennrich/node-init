GO111MODULE := on
DOCKER_TAG := $(or ${GITHUB_TAG_NAME}, latest)

all: node-init

.PHONY: node-init
node-init:
	go build -o bin/node-init
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

