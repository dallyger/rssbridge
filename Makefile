DOCKER_IMAGE ?= $(shell whoami)/$(shell basename `git rev-parse --show-toplevel`)

build: ## Build binaries
	go build -o ./bin/rssbridge ./cmd/rssbridge

do-build: ## Build docker image
	docker build -t dallyger/rssbridge:nightly .

do-run: ## Run docker image
	docker run --rm --name rssbridge -p 3000 dallyger/rssbridge:nightly

do-publish: ## Tag docker image with :latest and push it
	docker tag dallyger/rssbridge:nightly $(DOCKER_IMAGE):latest
	docker push $(DOCKER_IMAGE):latest
