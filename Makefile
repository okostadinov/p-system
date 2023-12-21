BINARY_NAME=p_system

prep:
	mkdir -p dist

build: prep
	go build -o ./dist/${BINARY_NAME} ./cmd/web/

compile: prep
	GOOS=linux GOARCH=amd64 go build -o ./dist/${BINARY_NAME}-linux ./cmd/web/
	GOOS=windows GOARCH=amd64 go build -o ./dist/${BINARY_NAME}-windows ./cmd/web/
	GOOS=darwin GOARCH=amd64 go build -o ./dist/${BINARY_NAME}-darwin ./cmd/web/

android: prep
	GOOS=android GOARCH=arm64 CGO_ENABLED=1 CC=$$ANDROID_NDK_HOME/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android34-clang go build -o ./dist/${BINARY_NAME}-android ./cmd/web/

run:
	go run ./cmd/web

clean:
	go clean
	rm -rf ./dist/
