BUILD_ARCHS = "linux/amd64,linux/arm64"
DOCKER_TAG = 3ayazaya/cscheck
VERSION = 1.0
APP_NAME = cscheck

build-and-push:
	docker buildx build --platform ${BUILD_ARCHS} --progress plain --pull -t ${DOCKER_TAG}:${VERSION} -t ${DOCKER_TAG}:latest --push -f Dockerfile . --no-cache

build-image:
	docker build -t ${DOCKER_TAG}:${VERSION} -t ${DOCKER_TAG}:latest . --no-cache

build:
	go build -o ${APP_NAME} cmd/cscheck/main.go