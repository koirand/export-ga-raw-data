BUILD_DIR := build
APP_NAME := export-ga-raw-data

# Install dependencies
deps:
	go get -v -d

clean:
	rm -rf $(BUILD_DIR)
	rm -f $(APP_NAME).zip

## Build binary
.PHONY: build
build: deps clean
	mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=386   go build -o $(BUILD_DIR)/windows_386.exe   main.go
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/windows_amd64.exe main.go
	cp credentials.json $(BUILD_DIR)
	zip -j ${APP_NAME}.zip $(BUILD_DIR)/*
