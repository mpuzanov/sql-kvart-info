.SILENT:
.PHONY: help clean build release run lint vet test
.DEFAULT_GOAL = build

APP=kvart-info
SOURCE=./cmd/${APP}
GOBASE=$(shell pwd)
RELEASE_DIR=$(GOBASE)/bin
NOW := $(shell date "+%Y-%m-%d %H-%M-%S")
Version=0.0.4

GO_SRC_DIRS := $(shell \
	find . -name "*.go" -not -path "./vendor/*" | \
	xargs -I {} dirname {}  | \
	uniq)

build: ## Build program
	$(call print-target)
	@go build -v -o ${APP} ${SOURCE}

run: build ## Run program
	$(call print-target)
	@./${APP}

test: ## go test with race detector and code covarage
	$(call print-target)
#	go test -v -cover ./internal/...
	gotestsum --format dots -- ./internal/... 
.PHONY: test

test-cover: ## go test with race detector and code covarage
	$(call print-target)
	go test -race -covermode=atomic -coverpkg=./... -coverprofile=c.out ./...
	go tool cover -func=c.out
#	go tool cover -html=c.out -o coverage.html
	@rm c.out
	
lint:  ## Lint the source files
	$(call print-target)
	@gofmt -s -w ${GO_SRC_DIRS}
	@go vet ${GO_SRC_DIRS}
	@golint ${GO_SRC_DIRS}

linter-golangci: ### check by golangci linter
	$(call print-target)
	golangci-lint run
.PHONY: linter-golangci	

clean: ## Clean build directory
	rm -f ./bin/${APP}
	rmdir ./bin
	rm -f coverage.*

mock:
	mockery --all
#	mockery --recursive=true --inpackage --name=repository
#	mockery --recursive=true --inpackage --name=usecase
.PHONY: mock

release: test lint
	$(call print-target)
	rm -rf ${RELEASE_DIR}${APP}*
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o ${APP} ${SOURCE}
	@GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o ${RELEASE_DIR}/${APP}.exe ${SOURCE}


define print-target
    @printf "Executing target: \033[36m$@\033[0m\n"
endef