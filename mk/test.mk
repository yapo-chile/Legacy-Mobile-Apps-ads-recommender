export REPORT_ARTIFACTS ?= reports

## Run tests and generate quality reports
test: build-test
	${DOCKER} run -ti --rm \
		-v $$(pwd)/${REPORT_ARTIFACTS}:/app/${REPORT_ARTIFACTS} \
		--env APPNAME \
		--env BRANCH \
		--name ${APPNAME}-test \
		${DOCKER_IMAGE}:test ${TEST_CMD}
	[[ "${TEST_CMD}" =~ coverhtml ]] && open ${REPORT_ARTIFACTS}/cover.html || true

## Build test docker image
build-test:
	${DOCKER} build \
		-t ${DOCKER_IMAGE}:test \
		-f docker/dockerfile.test \
		.

## Run tests and output coverage reports
cover: test-cover-int

## Run tests and open report on default web browser
coverhtml: test-coverhtml-int

## Run code linter and output report as text
checkstyle: test-checkstyle-int

## Run pact tests
pact: test-pact-int

.PHONY: test pact

# Internal targets are run on the test docker container,
# they are not intended to be run directly

cover-int:
	@scripts/commands/test_cover.sh cli

coverhtml-int:
	@scripts/commands/test_cover.sh html

checkstyle-int:
	@scripts/commands/test_style.sh display

pact-int:
	@scripts/commands/pact-test.sh

test-int:
	@echoHeader "Running Tests"
	@scripts/commands/test_style.sh
	@scripts/commands/test_cover.sh

test-%:
	make TEST_CMD="make $*" test
