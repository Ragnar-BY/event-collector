.PHONY: lint

lint:
	golangci-lint  run

run:
	go run ./...