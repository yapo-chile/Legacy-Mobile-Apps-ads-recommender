#!/bin/bash

set -e

mkdir -p ${REPORT_ARTIFACTS}

CHECKSTYLE_FILE=${REPORT_ARTIFACTS}/checkstyle-report.xml

echoHeader "Running Checkstyle Tests"

if [[ -n "$TRAVIS" ]]; then
    golangci-lint -c .golangci.yml run ./... | tee /dev/tty > ${CHECKSTYLE_FILE} && echo
else
    golangci-lint -c .golangci.yml --out-format "colored-line-number" run ./...
fi
status=${PIPESTATUS[0]}

# We need to catch error codes that are bigger then 2,
# they signal that gometalinter exited because of underlying error.
if [ ${status} -ge 1 ]; then
    echo "gometalinter exited with code ${status}, check gometalinter errors"
    exit ${status}
fi
exit 0
