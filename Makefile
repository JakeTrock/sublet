dev:
	go run ./*.go --debug
	cd coordinator && npm run dev
install:
	# Install Go dependencies
	go mod download
	go mod tidy
	
	cd coordinator && npm install

.PHONY: dev install
