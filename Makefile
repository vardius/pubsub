# This version-strategy uses git tags to set the version string
VERSION := $(shell git describe --tags --always --dirty)

# HELP
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help

help:
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

version: ## Show version
	@echo $(VERSION)

docker-build: ## Build given container. Example: `make docker-build`
	docker build -f Dockerfile --no-cache -t pubsub .

docker-run: ## Run container on given port. Example: `make docker-run PORT=9090`
	docker run -i -t --rm -p=$(PORT):$(PORT) --name="pubsub" pubsub

docker-stop: ## Stop docker container. Example: `make docker-stop`
	docker stop pubsub

docker-rm: docker-stop ## Stop and then remove docker container. Example: `make docker-rm`
	docker rm pubsub

docker-publish: docker-publish-latest docker-publish-version ## Docker publish. Example: `make docker-publish REGISTRY=https://your-registry.com`

docker-publish-latest: tag-latest
	@echo 'publish latest to $(REGISTRY)'
	docker push $(REGISTRY)/pubsub:latest

docker-publish-version: tag-version
	@echo 'publish $(VERSION) to $(REGISTRY)'
	docker push $(REGISTRY)/pubsub:$(VERSION)

docker-tag: docker-tag-latest docker-tag-version ## Tag current container. Example: `make docker-tag REGISTRY=https://your-registry.com`

docker-tag-latest:
	@echo 'create tag latest'
	docker tag pubsub $(REGISTRY)/pubsub:latest
	docker tag pubsub $(REGISTRY)/pubsub:latest

docker-tag-version:
	@echo 'create tag $(VERSION)'
	docker tag pubsub $(REGISTRY)/pubsub:$(VERSION)

docker-release: docker-build docker-publish ## Docker release - build, tag and push the container. Example: `make docker-release REGISTRY=https://your-registry.com`
