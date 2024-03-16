DOCKER_IMAGE ?= $(shell whoami)/$(shell basename `git rev-parse --show-toplevel`)

.PHONY: help
help: ## show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## build binaries
	go build -o ./bin/rssbridge ./cmd/rssbridge

.PHONY: do-build
do-build: build ## build docker image
	docker build -t dallyger/rssbridge:nightly .

.PHONY: do-run
do-run: ## run docker image
	docker run --rm --name rssbridge -p 3000:3000 dallyger/rssbridge:nightly

.PHONY: do-publish
do-publish: ## tag docker image with :latest and push it
	docker tag dallyger/rssbridge:nightly $(DOCKER_IMAGE):latest
	docker push $(DOCKER_IMAGE):latest
