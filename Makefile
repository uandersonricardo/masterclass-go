run:
	@bash ./scripts/run.sh

build:
	@go build -o bin/main cmd/main.go

compile-proto:
	@bash ./scripts/compile-proto.sh
