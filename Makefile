.PHONY: default
default: fmt lint install generate

.PHONY: testacc
testacc:
	TF_ACC=1 go test -v -cover ./...

.PHONY: fmt
fmt:
	gofmt -s -w -e .

.PHONY: lint
lint:
	golangci-lint run
