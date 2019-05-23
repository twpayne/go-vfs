.PHONY: nothing
nothing:

.PHONY: coverage.out
coverage.out:
	go test -coverprofile=$@ ./...

.PHONY: format
format:
	find . -name \*.go | xargs gofumports -w

.PHONY: html-coverage
html-coverage:
	go tool cover -html=coverage.out

.PHONY: install-tools
install-tools:
	GO111MODULE=off go get -u \
		github.com/golangci/golangci-lint/cmd/golangci-lint \
		github.com/mattn/goveralls \
		golang.org/x/tools/cmd/cover \
		mvdan.cc/gofumpt/gofumports

.PHONY: lint
lint:
	go vet ./...
	golangci-lint run

.PHONY: test
test:
	go test ./...