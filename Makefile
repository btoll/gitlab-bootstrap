CC=go
ARCH=amd64
MACOS=darwin
WIN=windows
BIN_NAME=gitlab-client
BIN_NATIVE=${BIN_NAME}_${GOARCH}-${GOOS}
BIN_MACOS=${BIN_NAME}_${ARCH}-${MACOS}
BIN_WIN=${BIN_NAME}_${ARCH}-${WIN}

.PHONY: build clean debug lint

build:
	${CC} build
	${CC} build -o ${BIN_NATIVE}
	GOARCH=${ARCH} GOOS=${MACOS} ${CC} build -o ${BIN_MACOS}
	GOARCH=${ARCH} GOOS=${WIN} ${CC} build -o ${BIN_WIN}

clean:
	${CC} clean
	rm -f ${BIN_NATIVE}
	rm -f ${BIN_MACOS}
	rm -f ${BIN_WIN}
	rm -f ${BIN_NAME}

debug:
	${CC} build
	dlv exec ./${BIN_NAME} -- --file examples/gitlab.json

lint:
	golangci-lint run

