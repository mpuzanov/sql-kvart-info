APP=kvart-info
SOURCE=main.go
GOBASE=$(shell pwd)
RELEASE_DIR=$(GOBASE)/bin

build: 
	@go build -v -o ${APP} ${SOURCE}

run:
	@go run ${SOURCE}

release:
	rm -rf ${RELEASE_DIR}${APP}*
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o ${APP} ${SOURCE}
	@GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o ${RELEASE_DIR}/${APP}.exe ${SOURCE}
