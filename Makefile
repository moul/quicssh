GO ?= go

.PHONY: install
install:
	$(GO) install .

.PHONY: test
test:
	$(GO) test -cover -v ./...

.PHONY: lint
lint:
	golangci-lint run --verbose ./...

.PHONY: release
release:
	goreleaser --snapshot --skip-publish --rm-dist
	@echo -n "Do you want to release? [y/N] " && read ans && [ $${ans:-N} = y ]
	goreleaser --rm-dist
