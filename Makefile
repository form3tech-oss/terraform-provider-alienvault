default: test

.PHONY: package
package:
	./scripts/package.sh

.PHONY: vet
vet:
	go vet ./...

.PHONY: vet
test:
	go test -v ./...

