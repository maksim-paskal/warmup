test:
	./scripts/validate-license.sh
	go fmt ./cmd/...
	go vet ./cmd/...
	go mod tidy
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run -v

build:
	go run github.com/goreleaser/goreleaser@latest build --rm-dist --snapshot
	mv ./dist/warmup_linux_amd64/warmup warmup
	docker build --pull . -t paskalmaksim/warmup:dev

push:
	docker push paskalmaksim/warmup:dev
