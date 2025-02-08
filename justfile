name := "rssbridge"
user := `whoami`

docker-image := env("DOCKER_IMAGE", user/name)

# show this help message
help:
    @just --list

# build binaries
build:
	go build -o ./bin/rssbridge .

# build image
[group("docker")]
do-build: build
	docker build \
		--label "org.opencontainers.image.created=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
		--label "org.opencontainers.image.authors=$(git config user.name)" \
		--label "org.opencontainers.image.source=$(git remote get-url origin)" \
		--label "org.opencontainers.image.version=$(git describe --tags --always --dirty)" \
		-t {{docker-image}}:nightly \
		.

# start container from image
[group("docker")]
do-run:
	docker run --rm --name rssbridge -p 3000:3000 rssbridge:nightly

# tag image with :latest and push it
[group("docker")]
do-publish:
	docker tag {{docker-image}}:nightly {{docker-image}}:latest
	docker push {{docker-image}}:latest
