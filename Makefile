.SILENT:

APP=kvart-info
SOURCE=main.go
GOBASE=$(shell pwd)
RELEASE_DIR=$(GOBASE)/bin

build:
	$(call print-target)
	@go build -v -o ${APP} ${SOURCE}

run:
	$(call print-target)
	@go run ${SOURCE}

release:
	$(call print-target)
	rm -rf ${RELEASE_DIR}${APP}*
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o ${APP} ${SOURCE}
	@GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o ${RELEASE_DIR}/${APP}.exe ${SOURCE}


define print-target
    @printf "Executing target: \033[36m$@\033[0m\n"
endef