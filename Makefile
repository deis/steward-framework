SHORT_NAME := steward-framework

REPO_PATH := github.com/deis/${SHORT_NAME}
DEV_ENV_IMAGE := quay.io/deis/go-dev:0.19.0
DEV_ENV_WORK_DIR := /go/src/${REPO_PATH}
DEV_ENV_PREFIX := docker run -it --rm -v ${CURDIR}:${DEV_ENV_WORK_DIR} -w ${DEV_ENV_WORK_DIR} -e K8S_CLAIMER_AUTH_TOKEN=${K8S_CLAIMER_AUTH_TOKEN}
DEV_ENV_CMD := ${DEV_ENV_PREFIX} ${DEV_ENV_IMAGE}

VERSION ?= "dev"
BINARY_DEST_DIR := rootfs/bin

all:
	@echo "Use a Makefile to control top-level building of the project."

# Allow developers to step into the containerized development environment
dev:
	${DEV_ENV_CMD} bash

bootstrap:
	${DEV_ENV_CMD} glide install

glideup:
	${DEV_ENV_CMD} glide up

test: test-unit

test-unit:
	${DEV_ENV_CMD} sh -c 'go test $$(glide nv)'

test-all:
	${DEV_ENV_CMD} sh -c 'go run testing/test_driver.go go test -tags integration $$(glide nv)'

test-cover:
	@${DEV_ENV_CMD} sh -c 'go run testing/test_driver.go _scripts/test-cover.sh'
