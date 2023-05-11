build:
	GOOS=darwin GOARCH=arm64 go build main.go
build-image:
	docker build -t rewriter .