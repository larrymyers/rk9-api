BINARY := rk9
BIN_PKG := larrymyers.com/rk9api/cmd/rk9
PKG_DIR := ./dist

TIMESTAMP := $(shell date -u "+%Y-%m-%dT%H:%M:%SZ")

ifdef GIT_COMMIT
COMMIT := ${GIT_COMMIT}
else
COMMIT := $(shell git rev-parse --short HEAD)
endif

ifdef GIT_BRANCH
BRANCH := ${GIT_BRANCH}
else
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
endif

LDFLAGS := -ldflags "-X main.Timestamp=${TIMESTAMP} -X main.Commit=${COMMIT} -X main.Branch=${BRANCH}"

all: clean test install

graphql:
	go run github.com/99designs/gqlgen generate

install:
	go install ${LDFLAGS} ${BIN_PKG}

run: install
	set -o allexport; source development.env; set +o allexport && rk9 -server

loader:
	set -o allexport; source development.env; set +o allexport && go run cmd/rk9-loader/main.go

dev_setup: install
	set -o allexport; source development.env; set +o allexport \
	&& rk9 -init-db

test:
	go test -v ./...

test_cov:
	go test -coverprofile=cover.out -v ./... && go tool cover -html=cover.out

lint:
	golint -set_exit_status ./...

vet:
	go vet ./...

build_debug:
	go build ${LDFLAGS} -gcflags "all=-N -l" -o rk9_debug ${BIN_PKG}

run_debug: build_debug
	dlv --listen=:2345 --headless=true --api-version=2 exec ./rk9_debug -- server

package:
	mkdir -p ${PKG_DIR}
	go build -o ${PKG_DIR}/${BINARY} ${LDFLAGS} ${BIN_PKG}

clean:
	rm -rf ${PKG_DIR}
	rm -rf ./test-results
	rm -f ${HOME}/go/bin/${BINARY}
	rm -f ./rk9-debug
	rm -f ./cover.out
