install:
	cd subletd && go mod tidy
	cd ts/scripts/ && go mod tidy
	cd ts && npm install --legacy-peer-deps

build-ts:
	cd ts && npm run build

startgo:
	(cd subletd && go run *.go)

startws:
	(cd ts && npm run dev)

start: startws startgo

nixbld:
	nix build

build-docker:
	docker build -f .devcontainer/Dockerfile .

clean:
	rm -rf ts/node_modules
	rm -rf ts/dist
	rm -rf ts/clients.db
	rm -rf go/client_id
	rm -rf result

# if you are running asdf you can run `asdf reshim go` to update the go path
# https://pkg.go.dev/gitlab.com/hmajid2301/optinix#section-readme
generate-options:
	go install gitlab.com/hmajid2301/optinix@latest
	optinix update
