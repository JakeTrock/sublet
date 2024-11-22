dev:
	./flatend/flatend --config ./flatend/config.toml & \
	go run ./*.go --debug & \
	cd coordinator && npm run dev
daemon:
	./flatend/flatend --config ./flatend/config.toml & go run ./*.go
install:
	# Install Go dependencies
	go mod download
	go mod tidy
	
	cd coordinator && npm install