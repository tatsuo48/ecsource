.PHONY: build install clean

build:
	@echo "Building..."
	@go build -o bin/ecsource main.go

install: build
	@echo "Installing..."
	@ln -sf $(PWD)/bin/ecsource /usr/local/bin/ecsource

clean:
	@echo "Cleaning..."
	@rm -rf /usr/local/bin/ecsource
