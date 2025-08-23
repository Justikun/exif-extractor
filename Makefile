build:
	@go build -o bin/ee

run: build
	@./bin/ee

test:
	@go test ./... -v

