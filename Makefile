default: fmt lint install generate

testacc:
	TF_ACC=1 go test -v -cover ./...

.PHONY: provider
provider:
	go run ./scripts/terraform-generator.go

.PHONY: docs
docs:
	go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@latest
	tfplugindocs generate

build:
	go build -v ./...

install: build
	go install -v ./...

fmt:
	gofmt -s -w -e .

lint:
	golangci-lint run

.PHONY: fmt lint testacc build install docs
