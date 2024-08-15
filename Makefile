build: clean
	@go build -o build/seedb cmd/main.go

run: build
	@./build/seedb

clean:
	@rm -rf ./build/