APP_NAME=lgn

CURRENT_VERSION=$(shell cat version)
NEXT_VERSION=$(shell ./increment_version.sh -m ${CURRENT_VERSION})

.PHONY: build-backend
build-backend:
	${MAKE} -C backend build

.PHONY: build-backend
run-backend:
	${MAKE} -C backend run

.PHONY: test-backend
test-backend:
	${MAKE} -C backend test

.PHONY: increment-version
increment-version:
	@echo ${NEXT_VERSION} > version
	@echo "New version: ${NEXT_VERSION}"

.PHONY: build-docker
build-docker:
	docker build -t koloooo/${APP_NAME}:${CURRENT_VERSION} -t koloooo/${APP_NAME}:latest .

.PHONY: push-docker
push-docker:
	docker push koloooo/${APP_NAME}:${CURRENT_VERSION}
	docker push koloooo/${APP_NAME}:latest
