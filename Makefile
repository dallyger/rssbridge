.PHONY=build
do-build: ## Build docker image
	docker build -t dallyger/rssbridge:nightly .
do-run: ## Run docker image
	docker run --rm --name rssbridge -p 3000 dallyger/rssbridge:nightly
