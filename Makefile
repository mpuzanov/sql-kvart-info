.SILENT:
.PHONY: help clean build release run lint vet test
.DEFAULT_GOAL = build

APP=kvart-info
SOURCE=main.go
GOBASE=$(shell pwd)
RELEASE_DIR=$(GOBASE)/bin

GO_SRC_DIRS := $(shell \
	find . -name "*.go" -not -path "./vendor/*" | \
	xargs -I {} dirname {}  | \
	uniq)
GO_TEST_DIRS := $(shell \
	find . -name "*_test.go" -not -path "./vendor/*" | \
	xargs -I {} dirname {}  | \
	uniq)

build:
	$(call print-target)
	@go build -v -o ${APP} ${SOURCE}

run: build
	$(call print-target)
	@./${APP}
#	@go run ${SOURCE}

lint:  ## Lint the source files
	$(call print-target)
	@gofmt -s -w ${GO_SRC_DIRS}
	@go vet ${GO_SRC_DIRS}
	@golint ${GO_SRC_DIRS}

release:
	$(call print-target)
	rm -rf ${RELEASE_DIR}${APP}*
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o ${APP} ${SOURCE}
	@GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o ${RELEASE_DIR}/${APP}.exe ${SOURCE}


define print-target
    @printf "Executing target: \033[36m$@\033[0m\n"
endef