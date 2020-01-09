default : vet test build

.PHONY: build
build:
	go build

.PHONY: vet
vet:
	go vet ./...

.PHONY: test
test:
	go test -v ./...

.PHONY: release
release:
	@curl -sL http://git.io/goreleaser | bash
