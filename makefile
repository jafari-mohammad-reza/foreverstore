test:
	@go test -v ./...
build:
	@go build -o dist
run: build
	@./dist