.PHONY: test
test:
	@./scripts/test.sh
build:
	go mod tidy
	@./scripts/validate-license.sh
	@./scripts/build-all.sh
	ls -lah _dist