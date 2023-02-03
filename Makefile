.PHONY: build

build:
	@echo "Building..."
	@go build -o bin/ecsresource main.go

install: build
	@echo "Installing..."
	@ln -sf $(PWD)/bin/ecsresource /usr/local/bin/ecsresource
