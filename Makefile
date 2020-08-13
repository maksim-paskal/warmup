.PHONY: test
test:
	@./scripts/test.sh
build:
	go mod tidy
	@./scripts/validate-license.sh
	docker build . -t paskalmaksim/warmup:dev