include mk/help.mk
include mk/colors.mk
SHELL=bash

# Image information
export APPNAME ?= $(shell basename `git rev-parse --show-toplevel`)
export BRANCH ?= $(shell git branch | sed -n 's/^\* //p')
export COMMIT ?= $(shell git rev-parse HEAD)
export COMMIT_DATE ?= $(shell TZ="America/Santiago" git show --quiet --date='format-local:%Y%m%d_%H%M%S' --format="%cd")
export COMMIT_DATE_UTC ?= $(shell TZ=UTC git show --quiet --date='format-local:%Y%m%d_%H%M%S' --format="%cd")
export CREATION_DATE ?= $(shell date -u '+%Y%m%d_%H%M%S')
export CREATOR ?= $(shell git log --format=format:%ae | head -n 1)

# Docker environment
export DOCKER_REGISTRY ?= registry.gitlab.com/yapo_team/legacy/mobile-apps
export DOCKER_IMAGE ?= ${DOCKER_REGISTRY}/${APPNAME}
export DOCKER_TAG ?= $(shell echo ${BRANCH} | tr '[:upper:]' '[:lower:]' | sed 's,/,_,g')
export DOCKER ?= docker

# Golang environment
export GO111MODULE ?= on

# K8s environment
export CHART_DIR ?= k8s/${APPNAME}

# Service variables
export SERVICE_PORT=8080
export BASE_URL="http://localhost:${SERVICE_PORT}"


## Setup a new service repository based on ads-recommender
clone:
	@scripts/commands/clone.sh

## Display basic service info
info:
	@echo "Service: ${APPNAME}"
	@echo "Images from latest commit:"
	@echo -e "- ${DOCKER_IMAGE}:${DOCKER_TAG}"
	@echo -e "- ${DOCKER_IMAGE}:${COMMIT_DATE_UTC}"
	@echo "API Base URL: ${BASE_URL}"
	@echo "Healthcheck: curl ${BASE_URL}/healthcheck"

include mk/dev.mk
include mk/test.mk
include mk/deploy.mk
