.PHONY: docker-install
docker-install:
	@CGO_ENABLED=0 go build -o ./app -ldflags "${LDFLAGS}" ${INPUT_NAME}

IMAGE_NAME=app-${APP_NAME}:${VERSION}

.PHONY: docker
docker:
	@if docker image inspect ${IMAGE_NAME} > /dev/null 2>&1; then docker rmi ${IMAGE_NAME}; fi
	@docker build -t ${IMAGE_NAME} .

.PHONY: image
image:
ifeq (${PLATFORM},${PLATFORM_MAC})
	@echo '[registry."${REPO}"]\nhttp = true' >${CURDIR}/build/docker/buildkitd.toml
	@if ! docker buildx ls | grep 'myBuilder' ; then docker buildx create --config=${CURDIR}/build/docker/buildkitd.toml --name myBuilder && docker buildx inspect --bootstrap; fi
	@docker buildx build --platform linux/amd64 -t ${REPO}/library/${IMAGE_NAME} . --push
else
ifeq (${PLATFORM},${PLATFORM_WIN})
	@echo "window build not support now..."
else
	@if docker image inspect ${IMAGE_NAME} > /dev/null 2>&1; then docker rmi ${IMAGE_NAME}; fi
	@docker build -t ${REPO}/library/${IMAGE_NAME} . --push
endif
endif
