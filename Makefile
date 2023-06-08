build:
	GOOS=darwin GOARCH=arm64 go build main.go
build-image:
	docker buildx build --platform linux/amd64,linux/arm/v7,linux/arm64/v8 -t ccr.ccs.tencentyun.com/k8s-test/rewriter:0.1 . --push