dev:
	go run ./*.go

install:
	# Install Go dependencies
	go mod download
	go mod tidy
	
	# Create Python virtual environment if it doesn't exist
	test -d coordinator/venv || python3 -m venv coordinator/venv
	
	# Install Python dependencies
	coordinator/venv/bin/pip install -r coordinator/requirements.txt

.PHONY: dev install
