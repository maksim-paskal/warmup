tag=dev
image=paskalmaksim/warmup:$(tag)

test:
	./scripts/validate-license.sh
	go fmt ./cmd/...
	go vet ./cmd/...
	go mod tidy
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run -v

build:
	go run github.com/goreleaser/goreleaser@latest build --rm-dist --snapshot
	mv ./dist/warmup_linux_amd64/warmup warmup
	docker build --pull . -t $(image)

push:
	docker push $(image)

scan:
	@trivy image \
	-ignore-unfixed --no-progress --severity HIGH,CRITICAL \
	$(image)