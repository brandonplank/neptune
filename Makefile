all: build

build:
	@echo "Building server for GNU/Linux"
	@echo "Building Web Assembly file"
	@GOOS=js GOARCH=wasm go build -o embed/public/assets/wasm/asm.wasm bin/webasm/main.go
	@echo "Building binary for the server"
	@go build -v
run:
	@GOOS=js GOARCH=wasm go build -o embed/public/assets/wasm/asm.wasm bin/webasm/main.go
	@go run .
