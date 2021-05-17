## Push docker image to containers.mpi-internal.com
docker-publish: build
	@echoTitle "Publishing docker image to Artifactory"
	${DOCKER} login --username "${ARTIFACTORY_USER}" --password "${ARTIFACTORY_PWD}" "${DOCKER_REGISTRY}"
	${DOCKER} push "${DOCKER_IMAGE}" --all-tags

## Upload helm charts for deploying on k8s
helm-publish:
	@echoHeader "Publishing helm package to Artifactory"
	helm lint ${CHART_DIR}
	helm package ${CHART_DIR}
	jfrog rt u "*.tgz" "helm-local/yapo/" || true

## Create production docker image
build:
	@echoHeader "Building production docker image"
	@set -x
	${DOCKER} build \
		-t ${DOCKER_IMAGE}:${DOCKER_TAG} \
		-f docker/dockerfile \
		--build-arg APPNAME=${APPNAME} \
		--build-arg GIT_COMMIT=${COMMIT} \
		--label appname=${APPNAME} \
		--label branch=${BRANCH} \
		--label build-date=${CREATION_DATE} \
		--label commit=${COMMIT} \
		--label commit-author=${CREATOR} \
		--label commit-date=${COMMIT_DATE} \
		.
	${DOCKER} tag ${DOCKER_IMAGE}:${DOCKER_TAG} ${DOCKER_IMAGE}:${COMMIT_DATE_UTC}
	@set +x

.PHONY: build
