default: fmt lint install generate

testacc:
	TF_ACC=1 go test -v -cover ./...

.PHONY: provider_code_spec.json
provider_code_spec.json:
	tfplugingen-openapi generate --config ./generator_config.yml --output ./provider_code_spec.json ./openapi/logto-openapi-source.json

.PHONY: provider
provider:
	tfplugingen-framework generate all --input ./provider_code_spec.json --output internal/provider

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