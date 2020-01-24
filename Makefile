default: test

.PHONY: package
package:
	./scripts/package.sh

.PHONY: vet
vet:
	go vet ./...

release:
	@curl -sL http://git.io/goreleaser | bash

.PHONY: vet
test:
	go test -v ./...

