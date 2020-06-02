GO=go

GOARCH=amd64
GOOS=linux

GOBUILD=${GO} build
GOTEST=${GO} test

TIMEOUT=-timeout 30s

COMMIT=$(shell git rev-parse --short=8 HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
TAG=${COMMIT}

BUILT=$(shell date -u '+%Y-%m-%d %H:%M:%S UTC')
LDFLAGS=-ldflags "-s -X 'github.com/vitpelekhaty/dbmill-cli/cmd/commands.GitBranch=${BRANCH}' -X 'github.com/vitpelekhaty/dbmill-cli/cmd/commands.GitCommit=${COMMIT}' -X 'github.com/vitpelekhaty/dbmill-cli/cmd/commands.Built=${BUILT}'"
BUILD_DIR=./bin

ifeq (${OS}, Windows_NT)
	FixPath = ${subst /,\,$1}
else
	FixPath = $1
endif

.PHONY: clean test build
.PHONY: test_internal_packages test_commands test_engine test_sqlserver_engine

all: build

clean:
	if [ -d "${BUILD_DIR}" ]; then rm -f "${BUILD_DIR}/*" ; else mkdir "${BUILD_DIR}" ; fi

test_dir_pkg:
	${GOTEST} ${TIMEOUT} github.com/vitpelekhaty/dbmill-cli/internal/pkg/dir

test_filter_pkg:
	${GOTEST} ${TIMEOUT} github.com/vitpelekhaty/dbmill-cli/internal/pkg/filter

test_internal_packages: test_dir_pkg test_filter_pkg

test_commands:
	${GOTEST} ${TIMEOUT} github.com/vitpelekhaty/dbmill-cli/cmd/commands

test_engine:
	${GOTEST} ${TIMEOUT} github.com/vitpelekhaty/dbmill-cli/cmd/engine

test_sqlserver_engine: test_engine
	${GOTEST} ${TIMEOUT} github.com/vitpelekhaty/dbmill-cli/cmd/engine/sqlserver

test: test_internal_packages test_commands test_engine test_sqlserver_engine

build: clean test
	GOOS=${GOOS} GOARCH=${GOARCH} ${GOBUILD} ${LDFLAGS} -o ${BUILD_DIR}/dbmill-cli .
