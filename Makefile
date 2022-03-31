CC=go
ARCH=amd64
MACOS=darwin
WIN=windows
BIN_NAME=gitlab-client
BIN_NATIVE=${BIN_NAME}_${GOARCH}-${GOOS}
BIN_MACOS=${BIN_NAME}_${ARCH}-${MACOS}
BIN_WIN=${BIN_NAME}_${ARCH}-${WIN}

.PHONY: build clean lint

build:
	${CC} build -o ${BIN_NATIVE}
	GOARCH=${ARCH} GOOS=${MACOS} ${CC} build -o ${BIN_MACOS}
	GOARCH=${ARCH} GOOS=${WIN} ${CC} build -o ${BIN_WIN}

clean:
	${CC} clean
	rm -f ${BIN_NATIVE}
	rm -f ${BIN_MACOS}
	rm -f ${BIN_WIN}

lint:
	golangci-lint run

